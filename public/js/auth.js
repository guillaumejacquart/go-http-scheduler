var auth = {
    getToken: function(){
        return localStorage.getItem('go-http-scheduler-token')
    },
    setToken: function(token){
        localStorage.setItem('go-http-scheduler-token', token)
    },
    logout: function(){
        localStorage.removeItem('go-http-scheduler-token')
        $('#login').modal('show');
    }
}

Vue.http.interceptors.push(function(request, next) {
    if (auth.getToken()) {
        request.headers.set('Authorization', 'Bearer ' + auth.getToken());
    }
    // continue to next interceptor
    next(function(response) {

        if (response.status == 401) {
            app.stopTimer();
            $('#login').modal('show');
        }

    });
});