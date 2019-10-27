import { ROOT_URL } from './root'
import axios from 'axios'

function getSingleSqlTrending(cluster, sqlId) {
  return axios.get(`${ROOT_URL}/api/v2/tiquerycnt?cluster=${cluster}&sqlID=${sqlId}`)
}

function getSqlInfo(cluster, queryId) {
  return axios.get(`${ROOT_URL}/api/v2/tiqueryinfo?cluster=${cluster}&queryID=${queryId}`)
}

export { getSingleSqlTrending, getSqlInfo }
