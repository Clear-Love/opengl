/*** 
 * @Author: lmio
 * @Date: 2023-02-26 17:23:48
 * @LastEditTime: 2023-02-26 17:26:20
 * @FilePath: /opengl/Teapot.cpp
 * @Description: 绘制茶壶
 */
// 引入OpenGL和GLUT库头文件
#include <stdlib.h>
#include <GL/glut.h>

// 初始化OpenGL的函数
void init(void)
{
    // 启用深度测试，避免遮挡问题 glEnable(GL_DEPTH_TEST);
    // 设置光源的位置
    GLfloat position[] = {1.0, 1.0, 1.0, 0.0};
    glLightfv(GL_LIGHT0, GL_POSITION, position);

    // 启用光照
    glEnable(GL_LIGHTING);
    // 启用第一个光源
    glEnable(GL_LIGHT0);

    // 定义材质的环境光分量，包括红、绿、蓝、透明度四个分量，取值范围为0到1
    GLfloat ambient[] = {0.0, 0.0, 0.0, 1.0};
    // 定义材质的散射光分量，包括红、绿、蓝、透明度四个分量，取值范围为0到1
    GLfloat diffuse[] = {0.5, 0.5, 0.5, 1.0};
    // 定义材质的镜面反射光分量，包括红、绿、蓝、透明度四个分量，取值范围为0到1
    GLfloat specular[] = {1.0, 1.0, 1.0, 1.0};
    // 将材质的环境光分量设置为定义的ambient值
    glMaterialfv(GL_FRONT, GL_AMBIENT, ambient);
    // 将材质的散射光分量设置为定义的diffuse值
    glMaterialfv(GL_FRONT, GL_DIFFUSE, diffuse);
    // 将材质的镜面反射光分量设置为定义的specular值
    glMaterialfv(GL_FRONT, GL_SPECULAR, specular);
    // 将材质的高光度设置为50.0，取值范围为0到128
    glMaterialf(GL_FRONT, GL_SHININESS, 50.0);
}

// 显示函数
void display(void)
{
    // 设置清除颜色
    glClearColor(0.75f, 0.75f, 0.75f, 1.0f);
    // 清除颜色缓冲区和深度缓冲区
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);

    // 创建一个显示列表
    glNewList(1, GL_COMPILE);
        // 在显示列表中创建一个茶壶对象，大小为0.5
        glutSolidTeapot(0.5);
    glEndList();

    // 调用显示列表1来渲染茶壶
    glCallList(1);

    // 刷新渲染结果
    glFlush();
}

// 窗口大小变化回调函数
void reshape(GLsizei w, GLsizei h)
{
    // 设置视口大小并进入投影模式
    glViewport(0, 0, w, h);
    glMatrixMode(GL_PROJECTION);
    glLoadIdentity();
    // 设置投影矩阵为正交投影，并设置视口的范围为-1到1的立方体
    glOrtho(-1.0, 1.0, -1.0, 1.0, -1.0, 1.0);
    // 返回模型视图模式
    glMatrixMode(GL_MODELVIEW);
}

// 主函数
int main(int argc, char** argv)
{
    // 初始化 GLUT 库
    glutInit(&argc, argv);

    // 设置显示模式，单缓冲、RGB 颜色模式、深度缓冲区
    glutInitDisplayMode(GLUT_SINGLE | GLUT_RGB | GLUT_DEPTH);

    // 设置窗口左上角在屏幕上的坐标位置
    glutInitWindowPosition(0, 0);

    // 设置窗口的大小
    glutInitWindowSize(500, 500);

    // 创建窗口并设置窗口标题
    glutCreateWindow(argv[0]);

    // 初始化状态
    init();

    // 注册窗口重置回调函数
    glutReshapeFunc(reshape);

    // 注册窗口绘制回调函数
    glutDisplayFunc(display);

    // 进入事件循环，等待用户操作
    glutMainLoop();

    return 0;
}