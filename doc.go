/*
hexa package is intended to support organizing services (parallel routines) with ease.
Currently only ContextStop() is recommened to use and others are experimental.

ContextStop is a imlementation of StoppableOne interface.
It also compatible to context.Context interface.
ContextStop is intended to ease offering StoppableOne interface.

ContextStop has two context internally. One is to receive external close request,
(by Close() method), another is to inform that the parallel routine has done.
*/
package hexa
