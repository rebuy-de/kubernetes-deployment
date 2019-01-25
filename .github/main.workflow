workflow "Test" {
  on = "push"
  resolves = ["Verify Golang Template"]
}

action "Verify Golang Template" {
  uses = "rebuy-de/golang-template@overhaul"
}
