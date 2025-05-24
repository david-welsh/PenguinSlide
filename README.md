# Penguin Slide

Penguin slide is a simple 2D game developed in Ebitengine and Go. It uses the 2D physics engine Chipmunk.
There isn't really a goal yet in the game, it has so far been developed only as an experiment in the usage of the above libraries.

## Playing the game

Run the game with
```shell
go run PenguinSlide
```

There are two worlds at present: 'Level 1' and 'Level 2'. Both have small level geometries to test out movement.

Controls:
Left and right arrow keys: Move left and right
Left shift: Start sliding
Space: Jump
Q and E: Control rotation left and right while in the air
Escape: Pause
F5: Restart level

F1: Open debug

## About the project

### Level system

The game supports multiple levels. Paths for levels are defined using SVGs (examples are in the `levels/` directory). These are converted for use by the game using the `tools/level-editor.html` web-page that features a simple Javascript script that converts the SVG curves to line segments and other level data (NB this process is very easy to do in JS using modern browsers advanced SVG support, but fairly difficult to achieve elsewhere hence the slightly circuitous route to create the level geometry). 
