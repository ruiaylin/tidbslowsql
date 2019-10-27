-- This Source Code Form is subject to the terms of the Mozilla Public
-- License, v. 2.0. If a copy of the MPL was not distributed with this
-- file, You can obtain one at http://mozilla.org/MPL/2.0/.
require "string"
local mysql = require "lpeg.mysql"
local tidb = require "lpeg.tidb"

function print_r ( t )
    local print_r_cache={}
    local function sub_print_r(t,indent)
        if (print_r_cache[tostring(t)]) then
            print(indent.."*"..tostring(t))
        else
            print_r_cache[tostring(t)]=true
            if (type(t)=="table") then
                for pos,val in pairs(t) do
                    if (type(val)=="table") then
                        print(indent.."["..pos.."] => "..tostring(t).." {")
                        sub_print_r(val,indent..string.rep(" ",string.len(pos)+8))
                        print(indent..string.rep(" ",string.len(pos)+6).."}")
                    elseif (type(val)=="string") then
                        print(indent.."["..pos..'] => "'..val..'"')
                    else
                        print(indent.."["..pos.."] => "..tostring(val))
                    end
                end
            else
                print(indent..tostring(t))
            end
        end
    end
    if (type(t)=="table") then
        print(tostring(t).." {")
        sub_print_r(t,"  ")
        print("}")
    else
        sub_print_r(t,"  ")
    end
    print()
end

local tidb_slow_query_log = {
[[
# Time: 2019-09-22T22:32:12.080095163+08:00
# Txn_start_ts: 411346595049111687
# User: reply_wr@10.27.75.75
# Conn_ID: 700288
# Query_time: 0.686471601
]],[[
# Time: 2019-08-14T09:26:59.487776265+08:00
# Txn_start_ts: 410450924122144769
# User: root@127.0.0.1
# Conn_ID: 3086
# Query_time: 1.527627037
# Process_time: 0.07 Request_count: 1 Total_keys: 131073 Process_keys: 131072 Prewrite_time: 0.335415029 Commit_time: 0.032175429 Get_commit_ts_time: 0.000177098 Local_latch_wait_time: 0.106869448 Write_keys: 131072 Write_size: 3538944 Prewrite_region: 1
# DB: test_1
# Is_internal: false
# Digest: 50a2e32d2abbd6c1764b1b7f2058d428ef2712b029282b776beb9506a365c0f1
# Stats: t:pseudo
# Num_cop_tasks: 1
# Cop_proc_avg: 0.07 Cop_proc_p90: 0.07 Cop_proc_max: 0.07 Cop_proc_addr: 172.16.5.87:20171
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr: 172.16.5.87:20171
# Mem_max: 525211
# Succ: false
insert into t select * from t;
]],[[
# Time: 2019-09-22T23:27:59.518031776+08:00
# Txn_start_ts: 0
# Query_time: 28.409603522
# Is_internal: true
# Digest: 2917a8c8bb18aff03c6dcb6b1afbbab3067e8d6934dc9f0817bf9ae95c543f1a
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
    insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 15, 0, 30537232),(237, 0, 2, 0, 1388056),(237, 0, 1, 0, 11104448),(237, 0, 9, 0, 4025528),(237, 0, 4, 0, 33660356),(237, 0, 16, 0, 11104448),(237, 0, 3, 0, 2082084),(237, 0, 6, 0, 1735070),(237, 0, 7, 0, 1735070),(237, 0, 5, 0, 1735070) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-09-23T21:48:57.889272146+08:00
# Txn_start_ts: 0
# User: dba@127.0.0.1
# Conn_ID: 831556
# Query_time: 2.211217815
# DB: information_schema
# Is_internal: false
# Digest: 9d8d1c27d43986bbe6bc047ca8b48e777b082b434033f32c3d9431740b31e35e
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
    insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 15, 0, 30537232),(237, 0, 2, 0, 1388056),(237, 0, 1, 0, 11104448),(237, 0, 9, 0, 4025528),(237, 0, 4, 0, 33660356),(237, 0, 16, 0, 11104448),(237, 0, 3, 0, 2082084),(237, 0, 6, 0, 1735070),(237, 0, 7, 0, 1735070),(237, 0, 5, 0, 1735070) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-09-22T22:51:20.44982795+08:00
# Txn_start_ts: 411346896187031569
# User: autopricewriter@10.27.2.32
# Conn_ID: 657128
# Query_time: 0.312873027
# Prewrite_time: 0.310802323 Commit_time: 0.001882215 Get_commit_ts_time: 8.9563e-05 Write_keys: 81 Write_size: 3829 Prewrite_region: 7
# DB: autoprice
# Is_internal: false
# Digest: 9505cacb7c710ed17125fcc6cb3669e8ddca6c8cd8af6a31f6b3cd64604c3098
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
commit;
]],[[
# Time: 2019-09-23T10:41:24.261911464+08:00
# Txn_start_ts: 411358064491364387
# User: reply_wr@10.27.74.26
# Conn_ID: 757559
# Query_time: 0.425277471
# Request_count: 1 Total_keys: 10 Process_keys: 9
# DB: Reply
# Index_ids: [1]
# Is_internal: false
# Digest: 76d4c8de80230dc532e28a55e193ddb0233be15dfb7d823de214ec36963e1e93
# Stats: ClubReply:411357980198961183
# Num_cop_tasks: 1
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
# Mem_max: 209
    insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 15, 0, 30537232),(237, 0, 2, 0, 1388056),(237, 0, 1, 0, 11104448),(237, 0, 9, 0, 4025528),(237, 0, 4, 0, 33660356),(237, 0, 16, 0, 11104448),(237, 0, 3, 0, 2082084),(237, 0, 6, 0, 1735070),(237, 0, 7, 0, 1735070),(237, 0, 5, 0, 1735070) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-09-23T11:13:10.432803591+08:00
# Txn_start_ts: 0
# User: dba@10.27.74.132
# Conn_ID: 780205
# Query_time: 1.7672986499999999
# DB: information_schema
# Is_internal: false
# Digest: b10085fa74b7b7012a8ca27ad260b7033b0284633a54c99747e45c68c3d16c5b
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
    insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 15, 0, 30537232),(237, 0, 2, 0, 1388056),(237, 0, 1, 0, 11104448),(237, 0, 9, 0, 4025528),(237, 0, 4, 0, 33660356),(237, 0, 16, 0, 11104448),(237, 0, 3, 0, 2082084),(237, 0, 6, 0, 1735070),(237, 0, 7, 0, 1735070),(237, 0, 5, 0, 1735070) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-09-23T11:13:47.67622078+08:00
# Txn_start_ts: 0
# User: dba@10.27.74.132
# Conn_ID: 780205
# Query_time: 2.084945073
# DB: information_schema
# Is_internal: false
# Digest: f1de62a09464e8a9f2c70b42ac315cbf8a6285222948e4b1b7b5dc02c77e492c
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0
    insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 15, 0, 30537232),(237, 0, 2, 0, 1388056),(237, 0, 1, 0, 11104448),(237, 0, 9, 0, 4025528),(237, 0, 4, 0, 33660356),(237, 0, 16, 0, 11104448),(237, 0, 3, 0, 2082084),(237, 0, 6, 0, 1735070),(237, 0, 7, 0, 1735070),(237, 0, 5, 0, 1735070) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-09-02T16:21:10.385337925+08:00
# Txn_start_ts: 410887770574946305
# Query_time: 15.483088376
# Prewrite_time: 11.174731535 Commit_time: 1.875035557 Get_commit_ts_time: 0.000378074 Local_latch_wait_time: 1.5926e-05 Write_keys: 16 Write_size: 453491 Prewrite_region: 1 Txn_retry: 1
# Is_internal: true
# Digest: 8cea8954b5a9d7e444feb61c1f5c240c0088d80c8cd4dfbb8a75a36a9579c799
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0 Cop_proc_addr:
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr:
# Succ: false
insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (81, 0, 31, 0, 12654317),(81, 0, 4, 0, 12977),(81, 0, 7, 0, 43162),(81, 0, 11, 0, 4073),(81, 0, 2, 0, 39566),(81, 0, 30, 0, 432031),(81, 0, 23, 0, 9320),(81, 0, 10, 0, 4115),(81, 0, 3, 0, 11386),(81, 0, 29, 0, 39555),(81, 0, 5, 0, 10529),(81, 0, 8, 0, 47520980),(81, 0, 6, 0, 17708),(81, 0, 9, 0, 12716454),(81, 0, 24, 0, 4660),(81, 0, 22, 0, 590) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-08-17T21:04:54.922748376+08:00
# Txn_start_ts: 410529848848285701
# Query_time: 3.036920735
# Prewrite_time: 0.058361101 Commit_time: 0.103531395 Get_commit_ts_time: 0.000149317 Local_latch_wait_time: 3.1903e-05 Write_keys: 10 Write_size: 308593 Prewrite_region: 2
# Is_internal: true
# Digest: 85aabb5a0425977cd016e15f3de5bb40e3bba1628c178ba32e386fefd4f42fa1
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0 Cop_proc_addr:
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr:
# Succ: false
insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 1, 0, 12416416),(237, 0, 15, 0, 34145144),(237, 0, 2, 0, 1552052),(237, 0, 9, 0, 4500292),(237, 0, 6, 0, 1940065),(237, 0, 3, 0, 2328078),(237, 0, 4, 0, 37637257),(237, 0, 16, 0, 12416416),(237, 0, 7, 0, 1940065),(237, 0, 5, 0, 1940065) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-08-17T21:04:02.863737213+08:00
# Txn_start_ts: 410529835911479297
# Query_time: 0.344850466
# Prewrite_time: 3.424110363 Commit_time: 0.000522605 Get_commit_ts_time: 0.0002014 Resolve_lock_time: 0.153877379 Local_latch_wait_time: 3.1218e-05 Write_keys: 10 Write_size: 308593 Prewrite_region: 2 Txn_retry: 2
# Is_internal: true
# Digest: 85aabb5a0425977cd016e15f3de5bb40e3bba1628c178ba32e386fefd4f42fa1
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0 Cop_proc_addr:
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr:
# Succ: false
insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 15, 0, 30537232),(237, 0, 2, 0, 1388056),(237, 0, 1, 0, 11104448),(237, 0, 9, 0, 4025528),(237, 0, 4, 0, 33660356),(237, 0, 16, 0, 11104448),(237, 0, 3, 0, 2082084),(237, 0, 6, 0, 1735070),(237, 0, 7, 0, 1735070),(237, 0, 5, 0, 1735070) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
]],[[
# Time: 2019-10-23T14:20:01.689536516+08:00
# Txn_start_ts: 412040980470169608
# User: gzconcern_wr@10.27.3.28
# Conn_ID: 408366
# Query_time: 0.10289279
# Local_latch_wait_time: 0.002043375 Write_keys: 4 Write_size: 178 Txn_retry: 1
# DB: gzconcern
# Is_internal: false
# Digest: 33310ef1dc065104278e13aa8d9da0a80f334b22f72026b378b818bfb10feb9d
# Num_cop_tasks: 0
# Cop_proc_avg: 0 Cop_proc_p90: 0 Cop_proc_max: 0 Cop_proc_addr: 
# Cop_wait_avg: 0 Cop_wait_p90: 0 Cop_wait_max: 0 Cop_wait_addr: 
# Succ: true
    insert into mysql.stats_histograms (table_id, is_index, hist_id, distinct_count, tot_col_size) values (237, 0, 15, 0, 30537232),(237, 0, 2, 0, 1388056),(237, 0, 1, 0, 11104448),(237, 0, 9, 0, 4025528),(237, 0, 4, 0, 33660356),(237, 0, 16, 0, 11104448),(237, 0, 3, 0, 2082084),(237, 0, 6, 0, 1735070),(237, 0, 7, 0, 1735070),(237, 0, 5, 0, 1735070) on duplicate key update tot_col_size = tot_col_size + values(tot_col_size);
 ]]
}

local function tidb_test1()
    print(tidb_slow_query_log[1])
    local t = tidb.tidb_slow_sql:match(tidb_slow_query_log[1])
    if not t then return error("no match") end
    print("data:", t)
    print_r( t )
end

local function tidb_test2(i)
    print(tidb_slow_query_log[i])
    local t = tidb.tidb21_slow_sql_grammar:match(tidb_slow_query_log[i])
    if not t then return error("no match") end
    print("data:", t)
    print_r( t )
end

--tidb_test1()
for i = 1 , #tidb_slow_query_log do
   print("---------------------------------------------------------")
   tidb_test2(i)
end



