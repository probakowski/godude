extern crate zstd;

use zstd::Encoder;

pub fn compress(data: &mut [u8]) -> Result<Vec<u8>, qoi::Error> {
    let mut out = Vec::new();
    let mut enc = Encoder::new(&mut out, 3)?;
    let enc2 = qoi::Encoder::new(data, 1, 1)?;
    enc2.encode_to_stream(&mut enc)?;
    enc.finish()?;
    Ok(out)
}
