package main

import(
  "bufio"
  "os"
)

func main() {
  stdin := bufio.NewScanner(os.Stdin)
  stdin.Scan()
  text := stdin.Text()
}
