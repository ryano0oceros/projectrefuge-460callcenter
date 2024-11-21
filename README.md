# Call Center Simulation

This project simulates the operations of a call center using discrete event simulation. The code benchmarks various combinations of parameters to measure the performance and efficiency of the call center, including metrics like total calls, abandoned calls, and agent utilization.

## Description

This simulation models the behavior of a call center with:
- A fixed number of agents handling incoming calls.
- A queue for calls when no agents are available.
- Event handling for call arrivals, completions, and abandonments.

The goal is to assess the performance of the call center under various conditions using metrics like:
- **Total calls handled**
- **Abandoned calls**
- **Agent utilization**

The simulation runs for a fixed period (`1440` minutes, simulating 8 hours) with configurable parameters.

## Input Parameters

The simulation accepts the following configurable input parameters:

1. **Number of Agents** (`numAgents`):
   - Range: 1–10
   - Description: The number of agents available to handle calls.

2. **Simulation Time** (`simulationTime`):
   - Fixed: 1440 minutes
   - Description: The duration of the simulation, representing an 8-hour day.

3. **Maximum Wait Time** (`maxWaitTime`):
   - Values: 5, 10, or 15 minutes
   - Description: The maximum time a call can wait in the queue before being abandoned.

4. **Lambda** (`lambda`):
   - Values: 0.5, 1.0, 1.5, or 2.0
   - Description: The arrival rate of calls, with higher values representing more frequent calls.

5. **Average Call Time** (`averageCallTime`):
   - Values: 3, 5, 7, or 9 minutes
   - Description: The average time an agent spends handling a call.

## Usage

### Running the Simulation

1. Clone the Repository:
```bash
git clone <repository-url>
cd [repo]
```

2. Run the Simulation: Ensure Go is installed on your system. Then run the simulation using:

```bash
go run main.go
```

This will run the simulation with all combinations of the input parameters and save the results.

## Output

The simulation generates a CSV file (simulation_results.csv) with the following columns:

- NumAgents: The number of agents in the simulation.
- SimulationTime: The total time for the simulation in minutes.
- MaxWaitTime: The maximum time a call can wait before abandonment.
- Lambda: The call arrival rate.
- AverageCallTime: The average duration of calls.
- TotalCalls: The total number of calls received.
- AbandonedCalls: The number of calls that were abandoned.
- Utilization: The percentage of time agents were busy during the simulation.

In the sample with existing input parameters, 480 simulations are run to cover all possible combinations of the input parameters.

## Excel Sorting for Analysis

You can open the CSV file in Excel to perform custom sorting for analysis. To identify a system configuration that minimizes call abandonment while maintaining high agent utilization, perform the following sort:

1. Primary Sort: AbandonedCalls (Smallest to Largest).
2. Secondary Sort: Utilization (Largest to Smallest).

This will help you find configurations that drop the fewest calls while keeping operators efficiently utilized.

## License
This project is licensed under the MIT License.
