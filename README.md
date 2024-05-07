# Concurrent Resource Access Simulation

This Go program simulates concurrent read and write operations on a shared resource using goroutines and context for timeout handling.

## Overview

The simulation consists of the following components:

1. **Resource**: Represents a shared resource that can be read from or written to. It includes methods for reading and writing data with read-write mutex for synchronization.

2. **Worker**: Represents a worker that performs read or write operations on the resource.

3. **Simulation Runner**: Runs the authentication server simulation with a specified number of workers and timeout duration. It creates a pool of workers, sets a timeout context, and simulates concurrent read and write operations.

## Installation

Ensure you have Go installed on your system. Then, clone the repository:

