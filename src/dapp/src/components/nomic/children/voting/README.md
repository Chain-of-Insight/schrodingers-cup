# Voting Child Component

## Requirements

- Send / receive votes in the UI
- Votes can be collected pooled on the dapp server  and sent to the API server as a batch containing all the player's votes
- The API or dapp server will need to track voting per session, because, by default, if a player downvotes a rule that is successfully passed they're awarded points for their dissent

### Turn Requirements

- The voting child component will also be used to determine what player's turn it is for voting yes or not to rule
- Turns must be timed - e.g. we can't let too many players hog the game window by dragging out turns. Since the overall game window is on a timer, we want to avoid players running up the clock or simply never voting until seconds before the window will close.
- If you miss your turn to vote you miss it, if you don't get to play another term because there's a lot of people in the game lobby that'll be a nice lesson to pay attention :D
