# rpi3/GUI/src/components/main/
Main component which is composed by the main panel (reservations and classrooms view) and the lateral panel (classroom status and occupation counter).

## img/
Images and figures used in the main panel.

## main.js
React implementation of the main component. It asks rpi3_api for data about reservations, the classrooms and their computers. It defines an integer global state that coordines what view must be showed every time.

## main.module.css
Styles applied to the main component and its panels.