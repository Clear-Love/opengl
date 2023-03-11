#include <iostream>
#include <vector>
#include <GL/glut.h>

using namespace std;

// 计算B样条基函数
double bspline(double t, int i, int p, vector<double> &U)
{
    if (p == 0)
    {
        if (U[i] <= t && t < U[i + 1])
            return 1;
        else
            return 0;
    }
    else
    {
        double a = (t - U[i]) / (U[i + p] - U[i]);
        double b = (U[i + p + 1] - t) / (U[i + p + 1] - U[i + 1]);
        return a * bspline(t, i, p - 1, U) + b * bspline(t, i + 1, p - 1, U);
    }
}

// 计算双三次B样条曲面
vector<vector<double>> bspline_surface(vector<vector<vector<double>>> &ctrl_pts,
                                       vector<double> &U, vector<double> &V,
                                       int u_num, int v_num, int p)
{
    double u_step = (U.back() - U.front()) / (u_num - 1);
    double v_step = (V.back() - V.front()) / (v_num - 1);

    vector<vector<double>> surface_pts(v_num, vector<double>(u_num * 3));

    for (int i = 0; i < u_num; i++)
    {
        for (int j = 0; j < v_num; j++)
        {
            double u = U.front() + i * u_step;
            double v = V.front() + j * v_step;

            // 计算曲面上的点
            vector<double> pt(3, 0);
            for (int k = 0; k < ctrl_pts.size(); k++)
            {
                for (int l = 0; l < ctrl_pts[0].size(); l++)
                {
                    double b_u = bspline(u, k, p, U);
                    double b_v = bspline(v, l, p, V);
                    for (int m = 0; m < 3; m++)
                    {
                        pt[m] += ctrl_pts[k][l][m] * b_u * b_v;
                    }
                }
            }
            surface_pts[j][i * 3] = pt[0];
            surface_pts[j][i * 3 + 1] = pt[1];
            surface_pts[j][i * 3 + 2] = pt[2];
        }
    }

    return surface_pts;
}

// 绘制B样条曲面
void draw_bspline_surface(vector<vector<double>> &surface_pts,
                          int u_num, int v_num)
{
    glColor3f(1.0, 1.0, 1.0);

    glBegin(GL_TRIANGLES);
    for (int i = 0; i < u_num - 1; i++)
    {
        for (int j = 0; j < v_num - 1; j++)
        {
            // 计算三角形顶点坐标
            vector<vector<double>> triangle_pts = {
                {surface_pts[j][i*3], surface_pts[j][i*3+1], surface_pts[j][i*3+2]},
                {surface_pts[j+1][i*3], surface_pts[j+1][i*3+1], surface_pts[j+1][i*3+2]},
                {surface_pts[j][i*3+3], surface_pts[j[i*3+4], surface_pts[j][i*3+5]}
};
            // 绘制三角形
            for (int k = 0; k < 3; k++)
            {
                glVertex3f(triangle_pts[k][0], triangle_pts[k][1], triangle_pts[k][2]);
            }
            for (int k = 0; k < 3; k++)
            {
                glVertex3f(triangle_pts[(k + 1) % 3][0], triangle_pts[(k + 1) % 3][1], triangle_pts[(k + 1) % 3][2]);
            }
            triangle_pts = {
                {surface_pts[j][i * 3 + 3], surface_pts[j][i * 3 + 4], surface_pts[j][i * 3 + 5]},
                {surface_pts[j + 1][i * 3], surface_pts[j + 1][i * 3 + 1], surface_pts[j + 1][i * 3 + 2]},
                {surface_pts[j + 1][i * 3 + 3], surface_pts[j + 1][i * 3 + 4], surface_pts[j + 1][i * 3 + 5]}};
            for (int k = 0; k < 3; k++)
            {
                glVertex3f(triangle_pts[k][0], triangle_pts[k][1], triangle_pts[k][2]);
            }
            for (int k = 0; k < 3; k++)
            {
                glVertex3f(triangle_pts[(k + 1) % 3][0], triangle_pts[(k + 1) % 3][1], triangle_pts[(k + 1) % 3][2]);
            }
        }
    }
    glEnd();
}

// OpenGL绘制函数
void display()
{
    glClear(GL_COLOR_BUFFER_BIT | GL_DEPTH_BUFFER_BIT);
    // 设置视角
    glMatrixMode(GL_MODELVIEW);
    glLoadIdentity();
    gluLookAt(0, 0, 5, 0, 0, 0, 0, 1, 0);

    // 绘制坐标轴
    glLineWidth(2.0);
    glBegin(GL_LINES);
    glColor3f(1.0, 0.0, 0.0);
    glVertex3f(0.0, 0.0, 0.0);
    glVertex3f(1.0, 0.0, 0.0);
    glColor3f(0.0, 1.0, 0.0);
    glVertex3f(0.0, 0.0, 0.0);
    glVertex3f(0.0, 1.0, 0.0);
    glColor3f(0.0, 0.0, 1.0);
    glVertex3f(0.0, 0.0, 0.0);
    glVertex3f(0.0, 0.0, 1.0);
    glEnd();

    // 设置曲面控制点
    vector<vector<vector<double>>> ctrl_pts = {
        {{-1, -1, 0}, {-1, 0, 2}, {-1, 1, 0}},
        {{0, -1, -2}, {0, 0, 4}, {0, 1, -2}},
        {{1, -1, 0}, {1, 0, 2}, {1, 1, 0}}};
    // 设置节点序列
    vector<double> U = {-1, -1, -1, 0, 1, 1, 1};
    vector<double> V = { -1, -1, -1, 0, 1, 1, 1};
    // 设置采样点密度
    int density = 50;

    // 生成曲面网格点
    vector<vector<double>> surface_pts = generate_surface(ctrl_pts, U, V, density);

    // 绘制曲面
    glColor3f(1.0, 1.0, 1.0);
    draw_surface(surface_pts);

    glutSwapBuffers();
}

int main(int argc, char **argv)
{
    // 初始化OpenGL
    glutInit(&argc, argv);
    glutInitDisplayMode(GLUT_RGB | GLUT_DOUBLE | GLUT_DEPTH);
    glutInitWindowSize(800, 600);
    glutCreateWindow("B-Spline Surface");
    // 设置OpenGL回调函数
    glutDisplayFunc(display);
    glEnable(GL_DEPTH_TEST);

    // 运行OpenGL主循环
    glutMainLoop();

    return 0;
}