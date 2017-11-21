/*
 SymnaTEC control - Measures muscle activity and controls a mechanical arm with it
 Copyright (c) Dorian Stoll 2017
 Licensed under the Terms of the MIT License
 */

package control

import (
    "gobot.io/x/gobot/platforms/raspi"
    "strconv"
)

const (
    HIGH = 1
    LOW  = 0
)


func MotorWorker(channel chan Measurement, MotorIN1 int, MotorIN2 int, MotorENA int, Threshold int) {

    // Start up the GPIO interface
    adapter := raspi.NewAdaptor()
    adapter.Connect()

    // Convert our integers into strings, since gobot uses strings
    IN1 := strconv.Itoa(MotorIN1)
    IN2 := strconv.Itoa(MotorIN2)
    ENA := strconv.Itoa(MotorENA)

    // Connect to the Motor Pins and set their signal to low
    err := adapter.DigitalWrite(IN1, LOW)
    if err != nil {
        panic(err)
    }
    err = adapter.DigitalWrite(IN2, LOW)
    if err != nil {
        panic(err)
    }
    err = adapter.PwmWrite(ENA, 0)
    if err != nil {
        panic(err)
    }

    // Iterate over the values in the channel
    for measurement := range channel {

        // Check if the extending muscle is active
        if measurement.Extending > 0 && measurement.Extending < Threshold {
            // Enable the IN1 channel
            adapter.DigitalWrite(IN1, HIGH)
        } else {
            // Disable the IN1 Channel
            adapter.DigitalWrite(IN1, LOW)
        }

        // Check if the flexing muscle is active
        if measurement.Flexing > 0 && measurement.Flexing < Threshold {
            // Enable the IN2 channel
            adapter.DigitalWrite(IN2, HIGH)
        } else {
            // Disable the IN2 Channel
            adapter.DigitalWrite(IN2, LOW)
        }

        // Set the speed value
        // We need to calculate some things, since the analog reading ranges from 0 to 1024, while the PWM
        // output only uses values from 0 to 255
        adapter.PwmWrite(ENA, byte(measurement.Speed / 1024 * 255))
    }

}