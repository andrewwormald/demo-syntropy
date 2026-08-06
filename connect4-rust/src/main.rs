mod board;
mod crlf;
mod game;
mod input;
mod renderer;

use std::io;
use std::os::fd::AsRawFd;
use std::process::ExitCode;

use crlf::CrlfWriter;

fn main() -> ExitCode {
    println!("Connect 4");

    let stdin = io::stdin();
    let fd = stdin.as_raw_fd();

    let original = match enable_raw_mode(fd) {
        Ok(original) => original,
        Err(err) => {
            eprintln!("connect4: {err}");
            return ExitCode::FAILURE;
        }
    };

    // Raw mode disables output post-processing, so a bare "\n" written
    // after this point no longer returns the cursor to column 0;
    // translate it to "\r\n" ourselves.
    let out = CrlfWriter::new(io::stdout());
    let result = game::run(stdin, out);

    disable_raw_mode(fd, &original);

    if let Err(err) = result {
        eprintln!("connect4: {err}");
        return ExitCode::FAILURE;
    }

    ExitCode::SUCCESS
}

/// Puts the terminal referenced by `fd` into raw mode (no line buffering,
/// no echo, and no signal-generating control characters, so Ctrl-C arrives
/// as the literal byte 0x03 for `input::decode` to handle as Quit) and
/// returns the original settings so they can be restored with
/// `disable_raw_mode`.
fn enable_raw_mode(fd: i32) -> io::Result<libc::termios> {
    unsafe {
        let mut original: libc::termios = std::mem::zeroed();
        if libc::tcgetattr(fd, &mut original) != 0 {
            return Err(io::Error::last_os_error());
        }

        let mut raw = original;
        libc::cfmakeraw(&mut raw);
        if libc::tcsetattr(fd, libc::TCSAFLUSH, &raw) != 0 {
            return Err(io::Error::last_os_error());
        }

        Ok(original)
    }
}

/// Restores terminal settings captured by `enable_raw_mode`.
fn disable_raw_mode(fd: i32, original: &libc::termios) {
    unsafe {
        libc::tcsetattr(fd, libc::TCSAFLUSH, original);
    }
}
