use crate::Hash;
use blake3;
use keccak;

#[derive(Clone, Copy)]
pub(super) struct PowHasher([u64; 10]);

#[derive(Clone, Copy)]
pub(super) struct HeavyHasher;

impl PowHasher {
    #[inline(always)]
    pub(super) fn new(pre_pow_hash: Hash, timestamp: u64) -> Self {
        let mut start = [0u64; 10];
        for (&pre_pow_word, num) in pre_pow_hash.0.iter().zip(start.iter_mut()) {
            *num = pre_pow_word;
        }
        start[4] = timestamp;
        Self(start)
    }

    #[inline(always)]
    pub(super) fn finalize_with_nonce(mut self, nonce: u64) -> Hash {
        self.0[9] = nonce;

        let mut input = [0u8; 80];

        for (i, chunk) in input.chunks_mut(8).enumerate() {
            chunk.copy_from_slice(self.0[i].to_le_bytes().as_slice());
        }

        let output = blake3::hash(&input);

        let mut output64 = [0u64; 4];

        for (i, chunk) in output.as_bytes().chunks(8).enumerate() {
            output64[i] = u64::from_le_bytes(chunk.try_into().unwrap());
        }

        Hash::new(output64)
    }
}

impl HeavyHasher {
    // The initial state of `cSHAKE256("ProofOfWorkHash")`
    // [4] -> 16654558671554924254 ^ 0x04(padding byte) = 16654558671554924250
    // [16] -> 9793466274154320918 ^ 0x8000000000000000(final padding) = 570094237299545110
    #[rustfmt::skip]
    const INITIAL_STATE: [u64; 25] = [
        4239941492252378377, 8746723911537738262, 8796936657246353646, 1272090201925444760, 16654558671554924250,
        8270816933120786537, 13907396207649043898, 6782861118970774626, 9239690602118867528, 11582319943599406348,
        17596056728278508070, 15212962468105129023, 7812475424661425213, 3370482334374859748, 5690099369266491460,
        8596393687355028144, 570094237299545110, 9119540418498120711, 16901969272480492857, 13372017233735502424,
        14372891883993151831, 5171152063242093102, 10573107899694386186, 6096431547456407061, 1592359455985097269,
    ];
    #[inline(always)]
    pub(super) fn hash(in_hash: Hash) -> Hash {
        let mut state = Self::INITIAL_STATE;
        for (&pre_pow_word, state_word) in in_hash.0.iter().zip(state.iter_mut()) {
            *state_word ^= pre_pow_word;
        }
        keccak::f1600(&mut state);
        Hash::new(state[..4].try_into().unwrap())
    }
}
