# Rules Child Component

## Requirements

- Finish and inherit the Practice module
- Connect with API callers and voting
- Must be able to call all CRUD API tasks for staging a rule for creation, update or removal
- Will only be staged for API / Tezos updates if voted in by the players or according to whatever other special considerations define how a rule can be added (e.g. if the players have changed the basic rules of the base voting system)
- Unlike votes we can stage rule changes immediately and have the API server pool them before running them all at the end of the game window

### Turn Requirements

- The rules child component will also be used to determine what player's turn it is
- Turns must be timed - e.g. we can't let too many players hog the game window by dragging out turns. Since the overall game window is on a timer, we want to avoid players running up the clock or simply never making a rules update request until seconds before the window will close.
- If you miss your turn you miss it, if you don't get to play another term because there's a lot of people in the game lobby that'll be a nice lesson to pay attention :D
