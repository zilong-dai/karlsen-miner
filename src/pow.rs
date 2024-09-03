use std::sync::Arc;

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

#[derive(Clone)]
pub struct State {
    // pub id: usize,
    matrix: Arc<Matrix>,
    pub target: Uint256,
    pub pow_hash_header: [u8; 72],
    // PRE_POW_HASH || TIME || 32 zero byte padding
    hasher: PowHasher,

    pub nonce_mask: u64,
    pub nonce_fixed: u64,
}

impl State {
    #[inline]
    pub fn new(block_seed: BlockSeed) -> Self {
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
        let matrix = Arc::new(Matrix::generate(pre_pow_hash));
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
        }
    }

    #[inline(always)]
    // PRE_POW_HASH || TIME || 32 zero byte padding || NONCE
    pub fn calculate_pow(&self, nonce: u64) -> Uint256 {
        // Hasher already contains PRE_POW_HASH || TIME || 32 zero byte padding; so only the NONCE is missing
        let hash = self.hasher.finalize_with_nonce(nonce);

        // v1
        self.matrix.heavy_hash(hash)

        // v2 todo
    }

    #[inline(always)]
    #[allow(dead_code)]
    pub fn check_pow(&self, nonce: u64) -> bool {
        let pow = self.calculate_pow(nonce);
        // The pow hash must be less or equal than the claimed target.
        pow <= self.target
    }
}
