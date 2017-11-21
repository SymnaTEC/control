/*
 SymnaTEC control - Measures muscle activity and controls a mechanical arm with it
 Copyright (c) Dorian Stoll 2017
 Licensed under the Terms of the MIT License
 */

package control

import (
    "github.com/SymnaTEC/go-adcpi"
    "time"
)

/*
 A type that stores the data that was measured from a muscle
 */
type Measurement struct {
    Flexing int
    Extending int
    Speed int
}

/**
 Reads the data from the muscle sensors in a defined interval. Must be run as a goroutine
 */
func MeasureWorker(channel chan Measurement, adc adcpi.Interface, flexingChannel byte, extendingChannel byte,
    speedChannel byte, speed int, interval float64) {

    // Measure data and write it to the channel
    // Passing true to the for loop creates an infinite loop that never stops, unless
    // The program is terminated.
    for true {

        // Read from the flexing muscle
        flexing := adc.ReadRaw(flexingChannel)

        // Read from the extending muscle
        extending := adc.ReadRaw(extendingChannel)

        // Read the speed value from a potentiometer, if dynamic speed is enabled
        if speed < 0 {
            speed = adc.ReadRaw(speedChannel)
        }

        // Send the values to the motor thread
        channel <- Measurement{ Flexing:flexing, Extending:extending, Speed:speed }

        // Wait for a certain amount of time
        time.Sleep(time.Duration(interval * 1000 * 1000 * 1000))
    }
}