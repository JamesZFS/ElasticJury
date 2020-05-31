module.exports = {
    "transpileDependencies": [
        "vuetify"
    ],
    devServer: {
        proxy: {
            '/api': {
                target: 'http://59.110.47.157/',
                changeOrigin: true,
                // pathRewrite: {
                //     '^/api': '/'
                // }
            },
        }
    }
}