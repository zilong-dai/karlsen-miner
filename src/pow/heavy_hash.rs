use crate::pow::{hasher::HeavyHasher, xoshiro::XoShiRo256PlusPlus};
use crate::Hash;
use std::mem::MaybeUninit;

#[derive(Debug, Ord, PartialOrd, Eq, PartialEq)]
pub struct Matrix(pub [[u16; 64]; 64]);

impl Matrix {
    #[inline(always)]
    pub fn generate(hash: Hash) -> Self {
        let mut generator = XoShiRo256PlusPlus::new(hash);
        loop {
            let mat = Self::rand_matrix_no_rank_check(&mut generator);
            if mat.compute_rank() == 64 {
                return mat;
            }
        }
    }

    #[inline(always)]
    fn rand_matrix_no_rank_check(generator: &mut XoShiRo256PlusPlus) -> Self {
        Self(array_from_fn(|_| {
            let mut val = 0;
            array_from_fn(|j| {
                let shift = j % 16;
                if shift == 0 {
                    val = generator.u64();
                }
                (val >> (4 * shift) & 0x0F) as u16
            })
        }))
    }

    #[inline(always)]
    fn convert_to_float(&self) -> [[f64; 64]; 64] {
        // SAFETY: An uninitialized MaybrUninit is always safe.
        let mut out: [[MaybeUninit<f64>; 64]; 64] = unsafe { MaybeUninit::uninit().assume_init() };

        out.iter_mut()
            .zip(self.0.iter())
            .for_each(|(out_row, mat_row)| {
                out_row
                    .iter_mut()
                    .zip(mat_row)
                    .for_each(|(out_element, &element)| {
                        out_element.write(f64::from(element));
                    })
            });
        // SAFETY: The loop above wrote into all indexes.
        unsafe { std::mem::transmute(out) }
    }

    pub fn compute_rank(&self) -> usize {
        const EPS: f64 = 1e-9;
        let mut mat_float = self.convert_to_float();
        let mut rank = 0;
        let mut row_selected = [false; 64];
        for i in 0..64 {
            if i >= 64 {
                // Required for optimization, See https://github.com/rust-lang/rust/issues/90794
                unreachable!()
            }
            let mut j = 0;
            while j < 64 {
                if !row_selected[j] && mat_float[j][i].abs() > EPS {
                    break;
                }
                j += 1;
            }
            if j != 64 {
                rank += 1;
                row_selected[j] = true;
                for p in (i + 1)..64 {
                    mat_float[j][p] /= mat_float[j][i];
                }
                for k in 0..64 {
                    if k != j && mat_float[k][i].abs() > EPS {
                        for p in (i + 1)..64 {
                            mat_float[k][p] -= mat_float[j][p] * mat_float[k][i];
                        }
                    }
                }
            }
        }
        rank
    }

    pub fn heavy_hash(&self, hash: Hash) -> Hash {
        let hash = hash.to_le_bytes();
        // SAFETY: An uninitialized MaybrUninit is always safe.
        let mut vec: [MaybeUninit<u8>; 64] = unsafe { MaybeUninit::uninit().assume_init() };
        for i in 0..32 {
            vec[2 * i].write(hash[i] >> 4);
            vec[2 * i + 1].write(hash[i] & 0x0F);
        }
        // SAFETY: The loop above wrote into all indexes.
        let vec: [u8; 64] = unsafe { std::mem::transmute(vec) };

        // Matrix-vector multiplication, convert to 4 bits, and then combine back to 8 bits.
        let mut product: [u8; 32] = array_from_fn(|i| {
            let mut sum1 = 0;
            let mut sum2 = 0;
            for (j, &elem) in vec.iter().enumerate() {
                sum1 += self.0[2 * i][j] * (elem as u16);
                sum2 += self.0[2 * i + 1][j] * (elem as u16);
            }
            ((sum1 >> 10) << 4) as u8 | (sum2 >> 10) as u8
        });

        // Concatenate 4 LSBs back to 8 bit xor with sum1
        product.iter_mut().zip(hash).for_each(|(p, h)| *p ^= h);
        HeavyHasher::hash(Hash::from_le_bytes(product))
    }
}

pub fn array_from_fn<F, T, const N: usize>(mut cb: F) -> [T; N]
where
    F: FnMut(usize) -> T,
{
    let mut idx = 0;
    [(); N].map(|_| {
        let res = cb(idx);
        idx += 1;
        res
    })
}
