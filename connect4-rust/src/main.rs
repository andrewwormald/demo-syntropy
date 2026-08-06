const ROWS: usize = 6;
const COLS: usize = 7;

fn main() {
    println!("Connect 4");
    print_placeholder_board();
}

fn print_placeholder_board() {
    for _ in 0..ROWS {
        let cells: Vec<&str> = vec!["."; COLS];
        println!("{}", cells.join(" "));
    }
}
