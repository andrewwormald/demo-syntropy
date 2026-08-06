const ROWS: usize = 6;
const COLS: usize = 7;

fn main() {
    println!("Connect 4");
    print_placeholder_board();
}

fn print_placeholder_board() {
    for row in placeholder_board_rows() {
        println!("{}", row);
    }
}

fn placeholder_board_rows() -> Vec<String> {
    (0..ROWS)
        .map(|_| vec!["."; COLS].join(" "))
        .collect()
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn placeholder_board_has_correct_dimensions() {
        let rows = placeholder_board_rows();
        assert_eq!(rows.len(), ROWS);
        for row in &rows {
            assert_eq!(row.split(' ').count(), COLS);
        }
    }

    #[test]
    fn placeholder_board_rows_are_all_dots() {
        let rows = placeholder_board_rows();
        for row in &rows {
            assert!(row.split(' ').all(|cell| cell == "."));
        }
    }
}
