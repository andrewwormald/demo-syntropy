//! An `io::Write` wrapper for talking to a terminal that has been put into
//! raw mode.

use std::io::{self, Write};

/// Wraps an `io::Write` and translates each "\n" not already preceded by
/// "\r" into "\r\n" before writing. Raw mode disables the tty driver's
/// output post-processing, so a lone "\n" no longer returns the cursor to
/// column 0; writes made while raw mode is active need to supply the "\r"
/// themselves.
pub struct CrlfWriter<W: Write> {
    w: W,
    last_byte: u8,
}

impl<W: Write> CrlfWriter<W> {
    /// Returns a `CrlfWriter` wrapping `w`.
    pub fn new(w: W) -> Self {
        CrlfWriter { w, last_byte: 0 }
    }
}

impl<W: Write> Write for CrlfWriter<W> {
    /// Translates "\n" to "\r\n" (unless already preceded by "\r",
    /// including across separate `write` calls) and writes the result to
    /// the underlying writer. Returns the number of bytes of `buf`
    /// consumed; on success this is always `buf.len()`, even though more
    /// bytes may have been written to the underlying writer.
    fn write(&mut self, buf: &[u8]) -> io::Result<usize> {
        let mut start = 0;
        let mut prev = self.last_byte;
        for (i, &b) in buf.iter().enumerate() {
            if b == b'\n' && prev != b'\r' {
                self.w.write_all(&buf[start..i])?;
                self.w.write_all(b"\r\n")?;
                start = i + 1;
            }
            prev = b;
        }
        if let Some(&last) = buf.last() {
            self.last_byte = last;
        }
        if start < buf.len() {
            self.w.write_all(&buf[start..])?;
        }
        Ok(buf.len())
    }

    fn flush(&mut self) -> io::Result<()> {
        self.w.flush()
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn translates_bare_newlines() {
        let mut buf = Vec::new();
        let mut w = CrlfWriter::new(&mut buf);

        let input = b"line one\nline two\nline three";
        let n = w.write(input).expect("write() should not error");

        assert_eq!(n, input.len());
        assert_eq!(buf, b"line one\r\nline two\r\nline three");
    }

    #[test]
    fn leaves_existing_crlf_alone() {
        let mut buf = Vec::new();
        let mut w = CrlfWriter::new(&mut buf);

        w.write_all(b"already\r\nfine\r\n").expect("write_all() should not error");

        assert_eq!(buf, b"already\r\nfine\r\n");
    }

    #[test]
    fn empty_input() {
        let mut buf = Vec::new();
        let mut w = CrlfWriter::new(&mut buf);

        let n = w.write(&[]).expect("write() should not error");

        assert_eq!(n, 0);
        assert!(buf.is_empty());
    }

    #[test]
    fn newline_split_across_writes_is_not_doubled() {
        let mut buf = Vec::new();
        let mut w = CrlfWriter::new(&mut buf);

        w.write_all(b"line one\r").expect("write_all() should not error");
        w.write_all(b"\nline two").expect("write_all() should not error");

        assert_eq!(buf, b"line one\r\nline two");
    }
}
