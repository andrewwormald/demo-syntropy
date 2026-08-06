mod board;
mod input;
mod renderer;

use board::Board;
use renderer::render;

fn main() {
    let b = Board::new();
    print!("{}", render(&b, 0));
}
