package main

import (
	"fmt"
	"libink/parser"
)


func main() {
	game := `room is "a palatial room"
room property furnishings can be empty "so very sad and empty"
room property furnishings can be kingly "like a king"
room property furnishings is empty
room can upgrade "" if "room property furnishings is empty" then "room change furnishings to kingly" 
throne is "a stately throne"
ball is "a beautiful blue bouncy ball"
player property status can be pauper "An ordinary guy"
player property status can be king "KING ME"
player can sing "" then "describe 'Your beautiful voice bounces off the walls like a ball'"
player can dance "" then "player get ball" "describe 'You are having a ball!'"
player can throw ball if "player has ball" then "room get throne" "describe 'The ball is throne!'" "room do upgrade ''"
player can sit throne if "room has throne" then "player change status to king"
room has player
win "You are king! YOU WIN!" if "player property status is king"
lose "That doesn't go with the decor at all! YOU LOSE!" if "room has ball"`
	u := parser.ParseGame(game)
	fmt.Println(u.Describe())
	u.Do("sing", "")
	u.Do("dance", "")
	u.Do("throw", "ball")
}
