# connect4-rust

Terminal ASCII Connect 4 for two players, sharing a keyboard.

## Build

```sh
cargo build
```

## Run

```sh
cargo run
```

## Controls

- Left / Right arrow keys — move the cursor between columns
- Enter — drop the current player's piece into the selected column

The board redraws after every drop. The game ends with an announcement
when a player connects four or the board fills up.

## Test

```sh
cargo test
```
