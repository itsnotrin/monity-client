# Monity Client
 
- [Monity Client](#monity-client)
- [What is monity?](#what-is-monity)
- [What is monity client for?](#what-is-monity-client-for)
- [What do I need to do here?](#what-do-i-need-to-do-here)
  - [What is monitorable?](#what-is-monitorable)
    - [What is planned to be monitorable?](#what-is-planned-to-be-monitorable)
      - [CPU:](#cpu)
      - [RAM:](#ram)
      - [GPU:](#gpu)
      - [Storage:](#storage)
      - [Ethernet:](#ethernet)

# What is monity?
Monity is an application in which you can aggregate a large chunk of data from all of your nodes- and monitor it all from one place!

# What is monity client for?
This is the 'client' that you run on your nodes. This is what the 'listener' server will make requests to in order to get the data it needs.
It essentially just reports back information that you have turned on for it to log.

# What do I need to do here?
Nothing- unless you want to help us develop further features, in which case, feel free to make a fork and make a PR or make an issue!
If you're looking for build instructions, you can find them in the [build.md file](https://github.com/itsnotrin/monity-client/blob/main/build.md)

## What is monitorable?
Essentially anything? If your distro can track it, you're welcome to make a PR and add it to the code, that way it's accessible!

### What is planned to be monitorable?

#### CPU:
- Individual Core Temperature
- Frequency (Clock speed) (Individual cores)
- CPU Usage
- Power Usage (?)
- CPU Fan Speed
- CPU Info (Brand, Name, SKU etc)

#### RAM:
- Speed (Frequency)
- Timings
- Temperature(?)
- Amount of RAM
- How much is being used
- How much SWAP/Page is being used.
- Name(?)

#### GPU:
- GPU Usage
- Fan Speed
- Core Clock
- Memory Clock
- Voltage
- Temperature
- Memory Usage
- Power Usage
- GPU Info (Brand, Name, SKU etc)

#### Storage:
- SMART Status
- Storage Usage
- Partition layouts
- Attached Drives
- Mount points
- Temperatures

#### Ethernet:
- Link Speed
- Link Status (Connected, Disconnected)
- IPv4 Address
- IPv6 Address (If enabled)
- Current Send / Receive speeds
- Latency to certain website of your chosing.