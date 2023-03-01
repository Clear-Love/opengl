/*** 
 * @Author: lmio
 * @Date: 2023-02-26 18:30:15
 * @LastEditTime: 2023-02-26 18:30:16
 * @FilePath: /opengl/bs.cpp
 * @Description: 绘制叶子
 */
#include <GL/glut.h>
#include<stdlib.h>
#include<stdio.h>
#include<math.h>
#include<time.h>
void fractalFern(void) {
    glClear(GL_COLOR_BUFFER_BIT);
    glColor3f(0.4f, 0.8f, 0.4f);
    glBegin(GL_POINTS);
    GLfloat x = 0, y = 0;
    srand((unsigned)time(0));

    for (int i = 0; i < 50000;i++) {
        double random = rand() / static_cast<double>(RAND_MAX);
        int n = (int)10000 * random;

        GLfloat _x, _y;
        if (n < 100) {
            _x = 0; _y = 0.16 * y;
        }
        else if (n < 800) {
            _x = 0.2 * x - 0.26 * y;
            _y = 0.23 * x + 0.22 * y + 1.6;
        }
        else if (n < 1500) {
            _x = -0.15 * x + 0.28 * y;
            _y = 0.26 * x + 0.24 * y + 0.44;
        }
        else {
            _x = 0.85 * x + 0.04 * y;
            _y = -0.04 * x + 0.85 * y + 1.6;
        }
        x = _x;
        y = _y;
        glVertex2f(_x / 10, _y / 10 - 0.3);
        glFlush();
    }
    glEnd();
    glFlush();
}

int main(int argc, char* argv[])
{
    glutInit(&argc, argv);
    glutInitDisplayMode(GLUT_RGB | GLUT_SINGLE);
    glutInitWindowPosition(100, 100);
    glutInitWindowSize(450, 450);
    glutCreateWindow("Barnsley");
    glutDisplayFunc(& fractalFern);
    glutMainLoop();
    return 0;
}
