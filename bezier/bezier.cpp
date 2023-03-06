/*** 
 * @Author: lmio
 * @Date: 2023-03-06 22:06:36
 * @LastEditTime: 2023-03-06 22:06:38
 * @FilePath: /opengl/bezier/bezier.cpp
 * @Description: 
 */
/*** 
 * @Author: lmio
 * @Date: 2023-03-06 22:00:13
 * @LastEditTime: 2023-03-06 22:00:15
 * @FilePath: /opengl-cpp/bezier.cpp
 * @Description: 
 */
#include <GL/glut.h>
#include <cmath>

// 控制点
GLfloat ctrlPoints[3][3] = {
    {-4.0, 0.0, 0.0},
    {0.0, 4.0, 0.0},
    {4.0, 0.0, 0.0}
};

// 绘制二次Bezier曲线
void drawQuadraticBezierCurve() {
    glMap1f(GL_MAP1_VERTEX_3, 0.0, 1.0, 3, 3, &ctrlPoints[0][0]);
    glEnable(GL_MAP1_VERTEX_3);

    glBegin(GL_LINE_STRIP);
    for (int i = 0; i <= 30; i++) {
        glEvalCoord1f((GLfloat) i / 30.0);
    }
    glEnd();
}

// 绘制场景
void display() {
    glClear(GL_COLOR_BUFFER_BIT);
    glColor3f(1.0, 1.0, 1.0);

    drawQuadraticBezierCurve();

    glFlush();
}

// 设置OpenGL状态
void init() {
    glClearColor(0.0, 0.0, 0.0, 0.0);
    glMatrixMode(GL_PROJECTION);
    glLoadIdentity();
    gluOrtho2D(-5.0, 5.0, -5.0, 5.0);
}

// 主函数
int main(int argc, char** argv) {
    glutInit(&argc, argv);
    glutInitDisplayMode(GLUT_SINGLE | GLUT_RGB);
    glutInitWindowSize(400, 400);
    glutInitWindowPosition(100, 100);
    glutCreateWindow("Quadratic Bezier Curve");
    init();
    glutDisplayFunc(display);
    glutMainLoop();
    return 0;
}