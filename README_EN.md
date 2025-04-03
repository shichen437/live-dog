<p align="center">
  <img src="web/src/assets/images/profile.jpg" height="128">
  <h1 align="center">Live-Dog</h1>
</p>

[简体中文](./README.md) ｜ English

# Intro

Live-Dog is a live broadcast recording tool built on GoFrame and RuoYi-Vue3 (currently only supports the Douyin/Bilibili platform).

*This project is only used for learning and communication and does not involve any private information (including storage, uploading, crawling, etc.).*

## Architecture

- Backend adoption [Goframe](https://github.com/gogf/gf)、MySQL
- Front-end adoption vue3、[RuoYi-Vue3](https://gitee.com/y_project/RuoYi-Vue)、[Element Plus](https://element-plus.org/zh-CN/)

## Built-in functions

1. System Management（Basic users, roles, menus, configuration management, etc.）
2. System Monitoring（Server, scheduled tasks, etc.）
3. Live broadcast management（Scheduled monitoring, cookies, broadcast push, etc.）

## System environment

- Golang: Golang1.20+
- Database: MySQL5.7+
- Node: 18+

## Rapid Development

1. Clone source code

    `git clone https://github.com/shichen437/live-dog`

2. Install yarn

    ```
    npm install -g yarn 
    ```

3. create database（eg：live-dog，recommended utf8mb4 character set）

4. Run server(Built-in database migrations)

    ```
    go run main.go
    ```

    If you are in a development environment, hot update can be used

    ```
    make cli.install #First run
    gf run main.go
    ```

5. Open web directory and start frontend

    ```
    yarn dev 
    ```

6. Login account

    Username: admin \
    Password: admin123

## Docker deployment

1. Pull image
    ```
    docker pull shichen437/live-dog:latest
    ```

2. Set environment variables
    <table>
    <tr align="center">
      <th>Variable</th>
      <th>Variable description</th>
      <th>Format</th>
      <th>Required</th>
    </tr>
    <tr align="center">
      <td>DATABASE_DEFAULT_LINK</td>
      <td>Database connection</td>
      <td>mysql:root:123456@tcp(192.168.3.16:13306)/live-dog?charset=utf8mb4&parseTime=true&loc=Local</td>
      <td>Yes</td>
    </tr>
    <tr align="center">
      <td>PROJECT_SM4KEY</td>
      <td>Sm4 encrypt key</td>
      <td>abcdefghijklmnopqrstuvwxyz123456 (32-length string)</td>
      <td>No</td>
    </tr>
    <tr align="center">
      <td>TZ</td>
      <td>Time zone</td>
      <td>Asia/Shanghai</td>
      <td>No</td>
    </tr>
    </table>

3. Run

    ```
    docker run -d --name live-dog --restart=always \
    -p 9876:9876 \
    -e DATABASE_DEFAULT_LINK="mysql:root:123456@tcp(127.0.0.1:3306)/db?charset=utf8mb4 \
    -e PROJECT_SM4KEY="abcdefghijklmnopqrstuvwxyz123456" \
    -e TZ="Asia/Shanghai" \
    -v your_video_path:/LiveDog/video \
    -v your_upload_path:/LiveDog/upload \
    shichen437/live-dog:latest
    ```

## Supported platforms

  <table>
    <tr align="center">
      <th>Platform</th>
      <th>Url Demo</th>
      <th>Cookie</th>
      <th>Remark</th>
    </tr>
    <tr align="center">
      <td>抖音</td>
      <td>https://live.douyin.com/123456789</td>
      <td>✅</td>
      <td></td>
    </tr>
  </table>

## License

  Mit

## Thanks

- Goframe <https://github.com/gogf/gf>
- golang-migrate <https://github.com/golang-migrate/migrate>
- RuoYi-Vue3 <https://gitee.com/y_project/RuoYi-Vue>
