import request from '@/utils/request'

// 查询角色列表
export function listAuthor(query) {
  return request({
    url: '/author/manage/list',
    method: 'get',
    params: query
  })
}

export function getAuthor(id) {
  return request({
    url: '/author/manage/' + id,
    method: 'get'
  })
}

// 新增角色
export function addAuthor(data) {
  return request({
    url: '/author/manage',
    method: 'post',
    data: data
  })
}

export function delAuthor(id) {
  return request({
    url: '/author/manage/' + id,
    method: 'delete',
  })
}

export function getTrend(query) {
  return request({
    url: '/author/manage/trend',
    method: 'get',
    params: query
  })
}