import { ROOT_URL } from './root'
import axios from 'axios'

function getRealtimeTrending(params) {
  const cluster = params.cluster
  const instance = params.instance
  const start = params.start
  const end = params.end

  return axios.get(
    `${ROOT_URL}/api/v2/tislowcnt?cluster=${cluster}&instName=${instance}&startTime=${start}&endTime=${end}`
  )
}

function getRealtimeQuery(params) {
  const cluster = params.cluster
  const instance = params.instance
  const start = params.start
  const end = params.end

  return axios.get(
    `${ROOT_URL}/api/v2/tislow?cluster=${cluster}&instName=${instance}&startTime=${start}&endTime=${end}`
  )
}

export { getRealtimeTrending, getRealtimeQuery }
