mod board;

const ROWS: usize = 6;
const COLS: usize = 7;

fn main() {
    for _ in 0..ROWS {
        for _ in 0..COLS {
            print!("| . ");
        }
        println!("|");
    }
    for _ in 0..COLS {
        print!("----");
    }
    println!("-");
}
