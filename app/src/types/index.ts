export interface User {
  ID: number
  username: string
  firstname: string
  lastname: string
  email: string
}

export interface JWTPayload {
  exp: number
  iat: number
  sub: string
}
