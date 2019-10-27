-- author  : linruichao(wechat:ruiaylin)
-- my email: 782082893@qq.com
-- company : autohome.com.cn
-- date    : 2019-09-28
-- func    : parse the tidb's slow query

l = require "lpeg"
local l = require "lpeg"
l.locale(l)
local tonumber = tonumber
local ip = require "lpeg.ip_address"

local M = {}
setfenv(1, M)

local unreserved    = l.alnum + l.S"-._~@[]:+,"
local null_reserved = l.alnum + l.S" -._~@[]:+,"
local space         = l.space^1
local sep           = l.P"\n"
local sql_end       = l.P";" * (l.P"\n" + -1)
local line          = (l.P(1) - l.P"\n")^0 * sep

local dot  = '.'
local exp  = l.S'eE'
local sign = l.S'+-'

local float         = sign^0 * l.digit^1 * l.P"."^0 * l.digit^0 * exp^0 * sign^0 * l.digit^0
local integer       = l.P"-"^-1 * l.digit^1 / tonumber
local true_false    = l.C(l.P"true" + "false")

local match = l.match

-- Time 
local time           = l.P"# Time: " * l.Cg(unreserved^0, "Time") * sep^0
-- Txn_start_ts
local txn_start_id   = "# Txn_start_ts: " * l.Cg(integer, "Txn_start_ts") * sep
-- User
local user           = "# User: " * l.Cg(unreserved^0, "User") * sep^0
-- Conn_ID
local conn_id        = "# Conn_ID: " * l.Cg(integer, "Conn_ID") * sep^0
-- Query_time
local query_time     = l.P"# Query_time: " * l.Cg(l.Ct(l.Cg(float , "value") * l.Cg(l.Cc"s", "representation")), "Query_time") * space^0 * sep^0

-- tikv
local tikv_start      = l.P"# "
local s_process_t     = l.P"# Process_time: " * l.Cg(l.Ct(l.Cg(float , "value") * l.Cg(l.Cc"s", "representation")), "Process_time")
local s_wait_t        = l.P"# Wait_time: " * l.Cg(l.Ct(l.Cg(float , "value") * l.Cg(l.Cc"s", "representation")), "Wait_time")
local s_backoff_t     = l.P"# Backoff_time: " * l.Cg(l.Ct(l.Cg(float , "value") * l.Cg(l.Cc"s", "representation")), "Backoff_time")

local process_t       = l.P"Process_time: " * l.Cg(l.Ct(l.Cg(float , "value") * l.Cg(l.Cc"s", "representation")), "Process_time")
local wait_t          = l.P"Wait_time: " * l.Cg(l.Ct(l.Cg(float , "value") * l.Cg(l.Cc"s", "representation")), "Wait_time")
local backoff_t       = l.P"Backoff_time: " * l.Cg(l.Ct(l.Cg(float , "value") * l.Cg(l.Cc"s", "representation")), "Backoff_time")

local s_request_count = l.P"# Request_count: " * l.Cg(integer, "Request_count")
local request_count   = l.P"Request_count: " * l.Cg(integer, "Request_count")
local total_keys      = l.P"Total_keys: " * l.Cg(integer, "Total_keys")
local process_keys    = l.P"Process_keys: " * l.Cg(integer, "Process_keys")

local s_prewrite_t    = l.P"# Prewrite_time: " * l.Cg(l.Ct(l.Cg(float, "value") * l.Cg(l.Cc"s", "representation")), "Prewrite_time")
local prewrite_t      = l.P"Prewrite_time: " * l.Cg(l.Ct(l.Cg(float, "value") * l.Cg(l.Cc"s", "representation")), "Prewrite_time")
local commit_t        = l.P"Commit_time: " * l.Cg(l.Ct(l.Cg(float, "value") * l.Cg(l.Cc"s", "representation")), "Commit_time")
local get_commit_ts_t = l.P"Get_commit_ts_time: " * l.Cg(l.Ct(l.Cg(float, "value") * l.Cg(l.Cc"s", "representation")), "Get_commit_ts_time")
local total_backoff_t = l.P"Total_backoff_time: " * l.Cg(l.Ct(l.Cg(float, "value") * l.Cg(l.Cc"s", "representation")), "Total_backoff_time")
local s_latch_wait_t  = l.P"# Local_latch_wait_time: " * l.Cg(l.Ct(l.Cg(float, "value") * l.Cg(l.Cc"s", "representation")), "Local_latch_wait_time")
local latch_wait_t    = l.P"Local_latch_wait_time: " * l.Cg(l.Ct(l.Cg(float, "value") * l.Cg(l.Cc"s", "representation")), "Local_latch_wait_time")
local resolve_lock_t  = l.P"Resolve_lock_time: " * l.Cg(l.Ct(l.Cg(float, "value") * l.Cg(l.Cc"s", "representation")), "Resolve_lock_time")

local write_keys      = l.P"Write_keys: " * l.Cg(integer, "Write_keys")
local write_size      = l.P"Write_size: " * l.Cg(integer, "Write_size")
local prewrite_region = l.P"Prewrite_region: " * l.Cg(integer, "Prewrite_region")
local txn_retry       = l.P"Txn_retry: " * l.Cg(integer, "Txn_retry")

local tikv_coprocesser= tikv_start * process_t^0 * space^0 * wait_t^0 * space^0 * backoff_t^0 * space^0 * request_count^0 * space^0 * total_keys^0 * space^0 * process_keys^0 * space^0 * prewrite_t^0 * space^0 * commit_t^0  * space^0 * get_commit_ts_t^0 * space^0 * resolve_lock_t^0 * space^0 * latch_wait_t^0 * space^0 * write_keys^0 * space^0 * write_size^0 * space^0 * prewrite_region^0 * space^0 * txn_retry^0 * space^0 * sep^0

local tikv_coprocesser01 = s_process_t * space^0 * wait_t^0 * space^0 * backoff_t^0 * space^0 * request_count^0 * space^0 * total_keys^0 * space^0 * process_keys^0 * space^0 * prewrite_t^0 * space^0 * commit_t^0 * space^0 * get_commit_ts_t^0 * space^0 * total_backoff_t^0 * space^0 * resolve_lock_t^0 * space^0 * latch_wait_t^0 * space^0 * write_keys^0 * space^0 * write_size^0 * space^0 * prewrite_region^0 * space^0 * txn_retry^0 * space^0 * sep^0
local tikv_coprocesser02 = s_wait_t * space^0 * wait_t^0 * space^0 * backoff_t^0 * space^0 * request_count^0 * space^0 * total_keys^0 * space^0 * process_keys^0 * space^0 * prewrite_t^0 * space^0 * commit_t^0 * space^0 * get_commit_ts_t^0 * space^0 * total_backoff_t^0 * space^0 * resolve_lock_t^0 * space^0 * latch_wait_t^0 * space^0 * write_keys^0 * space^0 * write_size^0 * space^0 * prewrite_region^0 * space^0 * txn_retry^0 * space^0 * sep^0
local tikv_coprocesser03 = s_backoff_t * space^0 * request_count^0 * space^0 * total_keys^0 * space^0 * process_keys^0 * space^0 * prewrite_t^0 * space^0 * commit_t^0 * space^0 * get_commit_ts_t^0 * space^0 * total_backoff_t^0 * space^0 * resolve_lock_t^0 * space^0 * latch_wait_t^0 * space^0 * write_keys^0 * space^0 * write_size^0 * space^0 * prewrite_region^0 * space^0 * txn_retry^0 * space^0 * sep^0
local tikv_coprocesser04 = s_request_count * space^0 * total_keys^0 * space^0 * process_keys^0 * space^0 * prewrite_t^0 * space^0 * commit_t^0 * space^0 * get_commit_ts_t^0 * space^0 * total_backoff_t^0 * space^0 * resolve_lock_t^0 * space^0 * latch_wait_t^0 * space^0 * write_keys^0 * space^0 * write_size^0 * space^0 * prewrite_region^0 * space^0 * txn_retry^0 * space^0 * sep^0
local tikv_coprocesser05 = s_prewrite_t * space^0 * commit_t^0 * space^0 * get_commit_ts_t^0 * space^0 * total_backoff_t^0 * space^0 * resolve_lock_t^0 * space^0 * latch_wait_t^0 * space^0 *  write_keys^0 * space^0 * write_size^0 * space^0 * prewrite_region^0 * space^0 * txn_retry^0 * space^0 * sep^0
local tikv_coprocesser06 = s_latch_wait_t * space^0 *  write_keys^0 * space^0 * write_size^0 * space^0 * prewrite_region^0 * space^0 * txn_retry^0 * space^0 * sep^0

local a_tikv_coprocess = (tikv_coprocesser01 + tikv_coprocesser02 + tikv_coprocesser03 + tikv_coprocesser04 + tikv_coprocesser05 + tikv_coprocesser06)
-- end testing

-- end testing
local schema_db   = l.P"# DB: " * l.Cg(unreserved^0, "DB") * sep^0 
--# Index_ids: [1]
local index_ids   = l.P"# Index_ids: " * l.Cg(unreserved^0, "Index_ids")  * sep^0
-- # Is_internal: false
local is_internal = l.P"# Is_internal: " * l.Cg(unreserved^0, "Is_internal") * sep^0
-- # Digest: fb0a8169a2b6dfd2c455d26335ef4238986bbd877db0cdd50b0c6c80dc9dc11f
local digest      = l.P"# Digest: " * l.Cg(unreserved^0, "Digest") * sep^0
 -- # Stats: ClubReply:411346498422833348
local stats       = l.P"# Stats: " * l.Cg(unreserved^0, "Stats") * sep^0 
-- cop_tasks 
local num_cop_tasks = l.P"# Num_cop_tasks: " * l.Cg(integer, "Num_cop_tasks") * sep^0

-- cop process
local proc_addr = l.P"Cop_proc_addr:" * space^0 * l.Cg(null_reserved^0, "Cop_proc_addr")
local cop_proc  = "# Cop_proc_avg: " * l.Cg(l.Ct(l.Cg(float / tonumber, "value")
                * l.Cg(l.Cc"s", "representation")), "Cop_proc_avg")
                * " Cop_proc_p90: " * l.Cg(l.Ct(l.Cg(float / tonumber, "value")
                * l.Cg(l.Cc"s", "representation")), "Cop_proc_p90")
                * " Cop_proc_max: " * l.Cg(l.Ct(l.Cg(float / tonumber, "value")
                * l.Cg(l.Cc"s", "representation")), "Cop_proc_max") * space^0 
                * proc_addr^0
                * sep^0

-- cop process
local wait_addr = l.P"Cop_wait_addr:" * l.Cg(null_reserved^0, "Cop_wait_addr") * l.P"\n"
local cop_wait  = "# Cop_wait_avg: " * l.Cg(l.Ct(l.Cg(float / tonumber, "value")
                  * l.Cg(l.Cc"s", "representation")), "Cop_wait_avg")
                  * " Cop_wait_p90: " * l.Cg(l.Ct(l.Cg(float / tonumber, "value")
                  * l.Cg(l.Cc"s", "representation")), "Cop_wait_p90")
                  * " Cop_wait_max: " * l.Cg(l.Ct(l.Cg(float / tonumber, "value")
                  * l.Cg(l.Cc"s", "representation")), "Cop_wait_max") * space^0 
                  * wait_addr^0 
                  
local mem_max   = "# Mem_max: " * l.Cg(integer, "Mem_max") * sep^0
-- # Succ: false
local succ      = "# Succ: " * l.Cg(true_false, "Succ") * sep^0
-- sql
local sql       = space^0 * l.Cg((l.P(1) - sql_end)^0 * sql_end, "Sql")

-- local first_seg  = time^0 * txn_start_id^0 * user^0 * conn_id^0 * query_time * tikv_coprocesser * space^0 * schema_db^0
local first_seg   = time^0 * txn_start_id^0 * user^0 * conn_id^0 * query_time * a_tikv_coprocess^0 * space^0 * schema_db^0
local second_seg  = first_seg * index_ids^0 * is_internal^0 * digest^0 * stats^0

tidb21_slow_sql_grammar = l.Ct(second_seg * num_cop_tasks^0 * space^0 * cop_proc^0 * space^0 * cop_wait^0 * space^0 * mem_max^0 * space^0 * succ^0 * space^0 *sql^0)

return M

