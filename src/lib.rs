pub mod pow;
pub mod target;

pub use pow::{State, BlockSeed, BlockVersion};
pub use target::Uint256;
type Hash = Uint256;

// use libc::{u32, u64, u8};

#[no_mangle]
pub extern "C" fn karlsen(
    work: *const [u64; 4],
    timestamp: u64,
    nonce: u64,
    res: *mut [u64; 4],
    log: u8,
) -> u32 {
    let ret = 0u32;
    let timestamp = timestamp as u64;
    let nonce = nonce as u64;
    //let (work, timestamp, nonce) = read_work(line.unwrap());
    let work: Uint256 = Uint256(unsafe { (*work as [u64; 4]).clone() });
    if log == 0 {
        println!("{:?},{:?},{:?}", work, timestamp, nonce);
    }

    let blookseed = pow::BlockSeed::PartialBlock {
        header_hash: work.0,
        timestamp,
        nonce,
        target: Default::default(),
        nonce_mask: Default::default(),
        nonce_fixed: Default::default(),
        hash: Default::default(),
    };
    let mut state = State::new(BlockVersion::V1, blookseed);
    let hash = state.calculate_pow(nonce);

    for i in 0..4 {
        unsafe { (*res)[i] = hash.0[i] };
    }
    ret
}
