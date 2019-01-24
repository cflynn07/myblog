var gulp = require('gulp')
var sass = require('gulp-sass')

// Compile sass into CSS
gulp.task('sass', function () {
  return gulp.src([
    'web/scss/*.scss'
  ]).pipe(sass())
    .pipe(gulp.dest('web/static/css'))
})

gulp.task('default', gulp.parallel(['sass']))
