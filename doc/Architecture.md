# Distributed File Search 

## System Architecture ##

* Golook Broker (see https://github.com/ottenwbe/golook):    
    * Features:            
        * Exposes an http api to query for files on different systems
        * Whitelisting: Users can control the files or folders that are exposed for queries by explicitly whitelisting them.  
    * Behaviour:
        * Brokers cache system and file information in a local repository
            * Goal: Reduce the time to query for a file
        * Monitors file changes and reports them to the cache
            * Only monitors files that have been explicitly reported by a user
    * Infrastructure:
        * Peer-to-peer infrastructure of brokers
        * Peer-to-peer infrastructure supports two scenarios:
            * Broadcast file reports, then query for files at any broker
            * Keep file reports locally, then broadcast queries to find files
    * Restrictions: 
        * All brokers have to be configured to the same scenario
  
* Golook CLI (see https://github.com/ottenwbe/golook-cli): 
    * Interface for users to report and query for reported files
    * Has an uplink connection to one broker which typically runs locally on the same host as the CLI.
    
## Software Architecture ##

* Four layers:
    * API Layer:
        * Controller
        * Repositories
    * Service Layer:
        * File Services
        * System Services
        * Configuration Service
    * Routing Layer:ma
        * Broadcast
    * Communication Layer:
        * RPC 
* Cross-cutting concerns:
    * Runtime:
        * Server
        * Commands