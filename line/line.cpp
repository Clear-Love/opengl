/*** 
 * @Author: lmio
 * @Date: 2023-03-06 22:15:24
 * @LastEditTime: 2023-03-06 22:15:25
 * @FilePath: /opengl/line/line.cpp
 * @Description: 
 */
/*** 
 * @Author: lmio
 * @Date: 2023-03-06 22:02:41
 * @LastEditTime: 2023-03-06 22:12:29
 * @FilePath: /opengl-cpp/line.cpp
 * @Description: 
 */
#include <GL/glut.h>

// 绘制直线
void drawLine() {
    glBegin(GL_LINES);
    glVertex2f(-0.5, 0.0);
    glVertex2f(0.5, 0.0);
    glEnd();
}

// 绘制场景
void display() {
    glClear(GL_COLOR_BUFFER_BIT);
    glColor3f(1.0, 1.0, 1.0);

    drawLine();

    glFlush();
}

// 设置OpenGL状态
void init() {
    glClearColor(0.0, 0.0, 0.0, 0.0);
    glMatrixMode(GL_PROJECTION);
    glLoadIdentity();
    gluOrtho2D(-1.0, 1.0, -1.0, 1.0);
}

// 主函数
int main(int argc, char** argv) {
    glutInit(&argc, argv);
    glutInitDisplayMode(GLUT_SINGLE | GLUT_RGB);
    glutInitWindowSize(400, 400);
    glutInitWindowPosition(100, 100);
    glutCreateWindow("Line");
    init();
    glutDisplayFunc(display);
    glutMainLoop();
    return 0;
}
