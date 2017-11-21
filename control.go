/*
 SymnaTEC control - Measures muscle activity and controls a mechanical arm with it
 Copyright (c) Dorian Stoll 2017
 Licensed under the Terms of the MIT License
 */

package main

import (
    "github.com/SymnaTEC/go-adcpi"
    "flag"
    "time"
    "github.com/SymnaTEC/control/control"
)

/*
 This is the entry point of the application. When the program is run, this will be the first function that gets called.
 It is responsible for creating the connection to the muscle sensor and starting the plotting tools.
 */
func main() {

    // Load the Settings
    LoadSettings()

    // Connect to the ADC Pi
    adc := adcpi.ADCPI(byte(Settings.Address), 12)

    // Make a channel to connect the thread that reads the muscle sensors and the one that controls the motor
    channel := make(chan control.Measurement)

    // Start the threads
    go control.MeasureWorker(channel, adc, byte(Settings.FlexingChannel), byte(Settings.ExtendingChannel),
        byte(Settings.SpeedChannel), Settings.Speed, Settings.Interval)
    go control.MotorWorker(channel, Settings.MotorIN1, Settings.MotorIN2, Settings.MotorENA, Settings.Threshold)


    // Leave the program running
    for true {
        time.Sleep(time.Second * 0.5)
    }
}

/*
 A type that stores all settings. These settings are loaded through command line arguments.
 */
type SettingsData struct {

    /*
     The I2C address of the interface we are connecting to. The default setting is 0x68 (so 104 in decimal notation).
     */
    Address int

    /*
     The channel of the analog pin where the muscle sensor for the flexing muscle is connected.
     */
    FlexingChannel int

    /*
     The channel of the analog pin where the muscle sensor for the extending muscle is connected.
     */
    ExtendingChannel int

    /*
     The GPIO Ports that control the motor
     */
    MotorIN1 int
    MotorIN2 int
    MotorENA int

    /*
     The amount of seconds that passes between two measurements
     */
    Interval float64

    /*
     The channel of the analog pin where the potentiometer for controlling the motor speed is connected.
     */
    SpeedChannel int

    /*
     The constant speed that the motor should use. If this is negative, the potentiometer is used.
     */
    Speed int

    /*
     If a measured value is below the threshold, the muscle is treated as active.
     */
    Threshold int
}

/*
 The Instance of the Settings Storage
 */
var Settings SettingsData

/*
 Parses the settings from the commandline parameters.
 */
func LoadSettings() {
    Settings = SettingsData{}
    flag.IntVar(&(Settings.Address), "address", 0x68, "The I2C address of the interface we " +
        "are connecting to.")
    flag.IntVar(&(Settings.FlexingChannel), "flexingchannel", 1, "The channel of the analog pin " +
        "where the muscle sensor for the flexing muscle is connected.")
    flag.IntVar(&(Settings.ExtendingChannel), "extendingchannel", 2, "The channel of the analog" +
        " pin where the muscle sensor for the extending muscle is connected.")
    flag.IntVar(&(Settings.MotorIN1), "motorin1", 38, "The GPIO Port where the IN1 channel for the" +
        " motor driver is connected")
    flag.IntVar(&(Settings.MotorIN2), "motorin2", 40, "The GPIO Port where the IN2 channel for the" +
        " motor driver is connected")
    flag.IntVar(&(Settings.MotorENA), "motorena", 35, "The GPIO Port where the ENA channel for the" +
        " motor driver is connected")
    flag.Float64Var(&(Settings.Interval), "interval", 0.1, "The amount of seconds that passes " +
        "between two measurements")
    flag.IntVar(&(Settings.SpeedChannel), "speedchannel", 3, "The channel of the analog pin " +
        "where the potentiometer for controlling the motor speed is connected.")
    flag.IntVar(&(Settings.Speed), "speed", -1, "The constant speed that the motor should use. " +
        "If this is negative, the potentiometer is used.")
    flag.IntVar(&(Settings.Threshold), "threshold", 100, "If a measured value is below the " +
        "threshold, the muscle is treated as active.")
    flag.Parse()
}