# TiDBSlowSQL
A Real time , High performance, Multi-dimension Slow Query Analysis System for TiDB 

## Modules

- Slow Query Collection module
  
  - gathering the slow query in stream
  - structure the text string to a structured data
  - send the data to kafka system
  
- SQL Message Processing module

  - consume messages from kafka
  - do sql parsing and gen the abstract SQL
  - computing the report and store it to MySQL

- The Web End

  - show the real time slow query stream 
  - show report data
