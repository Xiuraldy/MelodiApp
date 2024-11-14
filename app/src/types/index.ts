export interface User {
  ID: number
  username: string
  firstname: string
  lastname: string
  email: string
}

export interface Entry {
  ID: number
  title: string
  content: string
}

export interface JWTPayload {
  MapClaims: {
    eat: number
    iat: number
  }
  session: string
}
