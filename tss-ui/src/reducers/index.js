import * as types from '../actions/types'

const reducer = (
  state = {
    realtimeTrending: [],
    realtimeQuery: [],
    realtimeCluster: '',
    realtimeInstance: '',
    report: [],
    reportCluster: '',
    reportInstance: '',
    singleSqlTrending: [],
    sqlInfo: []
  },
  action
) => {
  switch (action.type) {
    case types.REALTIME_TRENDING_GET:
      return { ...state, realtimeTrending: action.data }
    case types.REALTIME_QUERY_GET:
      return { ...state, realtimeQuery: action.data }
    case types.REALTIME_CLUSTER_SET:
      return { ...state, realtimeCluster: action.cluster }
    case types.REALTIME_INSTANCE_SET:
      return { ...state, realtimeInstance: action.instance }
    case types.REPORT_GET:
      return { ...state, report: action.data }
    case types.REPORT_CLUSTER_SET:
      return { ...state, reportCluster: action.cluster }
    case types.REPORT_INSTANCE_SET:
      return { ...state, reportInstance: action.instance }
    case types.SINGLE_SQL_TRENDING_SET:
      return { ...state, singleSqlTrending: action.data }
    case types.SQL_INFO_SET:
      return { ...state, sqlInfo: Object.entries(action.sql) }
    default:
      return state
  }
}

export default reducer
