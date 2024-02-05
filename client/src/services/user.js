export function getAuthHeader(){
    if(token){
        return {Authorization: `Bearer ${getJWT()}`}
    }
    return {}
}

export function getJWT(){
    return localStorage.getItem('token')
}

export function setJWT(token){
    localStorage.setItem('token', token)
}

export function login({ email, password}){

}