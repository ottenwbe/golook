# Distributed File Search # 

## System Architecture ##

* Golook CLI (see https://github.com/ottenwbe/golook-cli): 
    * Interface for users to report and query for files
    * Uplink to one broker which typically runs locally on the same host as the cli
* Golook Broker (see https://github.com/ottenwbe/golook): 
    * Hierarchical Infrastructure
    * Caches file locations for downlink server
    * Handle queries for files
    * Client to further uplink servers
    
## Software Architecture ##