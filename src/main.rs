mod pow;
mod target;

use pow::{State, BlockVersion};
use target::Uint256;

use std::fs::File;
use std::io::{BufRead, BufReader};

type Hash = Uint256;

fn read_work(line: String) -> (Uint256, u64, u64) {
    let mut work = [0u64; 4];
    let mut nonce: u64 = 0;
    let mut timestamp: u64 = 0;
    let mut count = 0;

    // let pattern = line.split(',').map(f).collect();
    for str in line.split(',') {
        if count == 4 {
            timestamp = str.parse::<u64>().unwrap();
        }
        if count == 5 {
            let val = &str[2..];
            nonce = u64::from_str_radix(&val, 16).unwrap();
        }
        if count <= 3 {
            work[count] = str.parse::<u64>().unwrap();
        }
        count += 1;
    }
    (Uint256::new(work), timestamp, nonce)
}

fn main() {
    let pattern_file = "./testdata/pattern-v1.txt";
    let input = File::open(pattern_file).expect(&format!("open {} error", pattern_file));
    let buffered = BufReader::new(input);

    let target = Uint256::new([
        0xffffffffffffffff,
        0xffffffffffffffff,
        0xffffffffffffffff,
        0x0fffffff,
    ]);

    for line in buffered.lines() {
        // let work = [17940221284075783383u64, 5515539701051934179u64, 9306386394228168259u64, 13467475580129520626u64];
        // let timestamp = 1702373574430u64;
        // let nonce: u64 = 0x856072b92445a954;
        let linestr = line.unwrap();
        let (work, timestamp, nonce) = read_work(linestr);
        let blookseed = pow::BlockSeed::PartialBlock {
            header_hash: work.0,
            timestamp,
            nonce,
            target: target,
            nonce_mask: Default::default(),
            nonce_fixed: Default::default(),
            hash: Default::default(),
        };
        let mut state = State::new(BlockVersion::V1, blookseed);
        println!("result {:?}", state.check_pow(nonce));
    }
}
