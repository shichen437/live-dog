import request from '@/utils/request'

// 查询角色列表
export function listParseInfo(query) {
  return request({
    url: '/media/parse/list',
    method: 'get',
    params: query
  })
}

export function getParseInfo(id) {
  return request({
    url: '/media/parse/' + id,
    method: 'get'
  })
}

// 新增角色
export function parseUrl(data) {
  return request({
    url: '/media/parse',
    method: 'post',
    data: data
  })
}

export function delParseInfo(id) {
  return request({
    url: '/media/parse/' + id,
    method: 'delete'
  })
}

export function downloadMedia(data) {
  return request({
    url: '/media/parse/download',
    method: 'get',
    params: data
  })
}