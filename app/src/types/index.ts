export interface User {
  ID: number
  username: string
  firstname: string
  lastname: string
  email: string
}

export interface Person {
  age: number
  workclass: string
  fnlwgt: number
  education: string
  'education-num': number
  'marital-status': string
  occupation: string
  relationship: string
  race: string
  sex: string
  'capital-gain': number
  'capital-loss': number
  'hours-per-week': number
  'native-country': string
  income: string
}

export interface JWTPayload {
  MapClaims: {
    eat: number
    iat: number
  }
  session: string
}
