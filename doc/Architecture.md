# Distributed File Search # 

## System Architecture ##

* CLI: 
    * Interface for users. 
    * Uplink to one broker in order to report and query for file locations.
* Golook Broker: 
    * Hierarchical Infrastructure
    * Caches file locations for downlink server
    * Handle queries for files
    * Client for further uplink servers
    
## Software Architecture ##