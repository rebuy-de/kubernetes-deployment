workflow "Push" {
  on = "push"
  resolves = ["Verify Golang Template"]
}

action "Verify Golang Template" {
  uses = "rebuy-de/golang-template@v3.1.0"
}
