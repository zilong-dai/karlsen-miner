use fish_hash::{Context, Hash512, HashData};

use crate::{
    pow::{hasher::PowHasher, heavy_hash::Matrix},
    target::Uint256,
    Hash,
};

mod hasher;
mod heavy_hash;
mod xoshiro;

#[derive(Clone, Debug)]
pub enum BlockSeed {
    // FullBlock(Box<RpcBlock>),
    PartialBlock {
        header_hash: [u64; 4],
        timestamp: u64,
        nonce: u64,
        target: Uint256,
        nonce_mask: u64,
        nonce_fixed: u64,
        hash: Option<String>,
    },
}

#[derive(Clone, Debug)]
#[allow(dead_code)]
pub enum BlockVersion {
    V1,
    V2,
}

pub struct State {
    // pub id: usize,
    matrix: Option<Matrix>,
    pub target: Uint256,
    pub pow_hash_header: [u8; 72],
    // PRE_POW_HASH || TIME || 32 zero byte padding
    hasher: PowHasher,

    pub nonce_mask: u64,
    pub nonce_fixed: u64,
    pub context: Option<Context>,

    pub version: BlockVersion,
}

impl State {
    #[inline]
    pub fn new(version: BlockVersion, block_seed: BlockSeed) -> Self {
        let pre_pow_hash;
        let header_timestamp: u64;
        let header_target;
        let nonce_mask: u64;
        let nonce_fixed: u64;
        match block_seed {
            BlockSeed::PartialBlock {
                ref header_hash,
                ref timestamp,
                ref target,
                nonce_fixed: fixed,
                nonce_mask: mask,
                ..
            } => {
                pre_pow_hash = Hash::new(*header_hash);
                header_timestamp = *timestamp;
                header_target = *target;
                nonce_mask = mask;
                nonce_fixed = fixed
            }
        }

        // PRE_POW_HASH || TIME || 32 zero byte padding || NONCE
        let hasher = PowHasher::new(pre_pow_hash, header_timestamp);

        let context = match version {
            BlockVersion::V1 => None,
            BlockVersion::V2 => Some(Context::new(false, None)),
        };

        let matrix = match version {
            BlockVersion::V1 => Some(Matrix::generate(pre_pow_hash)),
            BlockVersion::V2 => None,
        };

        let mut pow_hash_header = [0u8; 72];

        pow_hash_header.copy_from_slice(
            [
                pre_pow_hash.to_le_bytes().as_slice(),
                header_timestamp.to_le_bytes().as_slice(),
                [0u8; 32].as_slice(),
            ]
            .concat()
            .as_slice(),
        );
        Self {
            matrix,
            target: header_target,
            pow_hash_header,
            hasher,
            nonce_mask,
            nonce_fixed,
            context,
            version,
        }
    }

    #[inline(always)]
    // PRE_POW_HASH || TIME || 32 zero byte padding || NONCE
    pub fn calculate_pow(&mut self, nonce: u64) -> Uint256 {
        // Hasher already contains PRE_POW_HASH || TIME || 32 zero byte padding; so only the NONCE is missing
        let hash = self.hasher.finalize_with_nonce(nonce);
        match self.version {
            BlockVersion::V1 => self
                .matrix
                .as_ref()
                .expect("matrix unwrap error")
                .heavy_hash(hash),
            BlockVersion::V2 => {
                let mut seed = [0u8; 64];

                seed[0..32].copy_from_slice(&hash.to_le_bytes());

                let mid_hash = fish_hash::fishhash_kernel(
                    &mut self.context.as_mut().expect("context unwrap error"),
                    &Hash512::new_from(seed),
                );

                let output = blake3::hash(&mid_hash.as_bytes());

                let mut output64 = [0u64; 4];

                for (i, chunk) in output.as_bytes().chunks(8).enumerate() {
                    output64[i] = u64::from_le_bytes(chunk.try_into().unwrap());
                }

                Uint256::new(output64)
            }
        }
    }

    #[inline(always)]
    #[allow(dead_code)]
    pub fn check_pow(&mut self, nonce: u64) -> bool {
        let pow = self.calculate_pow(nonce);
        println!("pow:");
        for num in pow.to_le_bytes(){
            print!("{:02x}", num);
        }
        println!();
        // The pow hash must be less or equal than the claimed target.
        pow <= self.target
    }
}
