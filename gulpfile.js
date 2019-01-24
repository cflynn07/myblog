var gulp = require('gulp')
var sass = require('gulp-sass')

// Compile sass into CSS
gulp.task('sass', function () {
  return gulp.src([
    'web/scss/*.scss'
  ]).pipe(sass())
    .pipe(gulp.dest('web/static/css'))
})

// configure which files to watch and what tasks to use on file changes
gulp.task('sass-watch', function() {
  return gulp.watch('web/scss/**/*.scss', gulp.series('sass'));
});

gulp.task('default', gulp.parallel(['sass']))
