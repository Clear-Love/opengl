#include <GL/glut.h>
#include <vector>
#include <iostream>

using namespace std;

// 控制点结构体
struct Point
{
    double x, y;
};

// B样条曲线类
class BSpline
{
public:
    // 构造函数
    BSpline() : m_k(4), m_t(0.0) {}

    // 添加控制点
    void addPoint(double x, double y)
    {
        Point p;
        p.x = x;
        p.y = y;
        m_points.push_back(p);
    }

    // 计算B样条曲线
    void compute()
    {
        // 计算节点向量
        int n = m_points.size() - 1;
        int m = n + m_k + 1;
        for (int i = 0; i <= m; i++)
        {
            if (i < m_k)
            {
                m_t.push_back(0.0);
            }
            else if (i > n)
            {
                m_t.push_back(1.0);
            }
            else
            {
                m_t.push_back((double)(i - m_k + 1) / (n - m_k + 2));
            }
        }

        // 计算曲线
        for (double u = 0.0; u <= 1.0; u += 0.01)
        {
            Point p = calculate(u);
            m_curve.push_back(p);
        }
    }

    // 绘制B样条曲线
    void draw()
    {
        glColor3f(1.0, 0.0, 0.0);
        glLineWidth(2.0);
        glBegin(GL_LINE_STRIP);
        for (int i = 0; i < m_curve.size(); i++)
        {
            glVertex2f(m_curve[i].x, m_curve[i].y);
        }
        glEnd();
        glColor3f(0.0, 0.0, 1.0);
        glPointSize(5.0);
        glBegin(GL_POINTS);
        for (int i = 0; i < m_points.size(); i++)
        {
            glVertex2f(m_points[i].x, m_points[i].y);
        }
        glEnd();
    }

private:
    // 计算B样条曲线上的点
    Point calculate(double u)
    {
        int n = m_points.size() - 1;
        int k = m_k;
        vector<double> t = m_t;

        double x = 0.0;
        double y = 0.0;

        for (int i = 0; i <= n; i++)
        {
            double b = basis(i, k, u, t);
            x += b * m_points[i].x;
            y += b * m_points[i].y;
        }

        Point p;
        p.x = x;
        p.y = y;
        return p;
    }

    // 计算B样条基函数
    double basis(int i, int k, double u, vector<double> t)
    {
        if (k == 1)
        {
            if (u >= t[i] && u < t[i + 1])
            {
                return 1.0;
            }
            else
            {
                return 0.0;
            }
            else
            {
                double w1 = 0.0;
                double w2 = 0.0;

                if (t[i + k - 1] != t[i])
                {
                    w1 = (u - t[i]) / (t[i + k - 1] - t[i]) * basis(i, k - 1, u, t);
                }

                if (t[i + k] != t[i + 1])
                {
                    w2 = (t[i + k] - u) / (t[i + k] - t[i + 1]) * basis(i + 1, k - 1, u, t);
                }

                return w1 + w2;
            }
        }

    private:
        vector<Point> m_points; // 控制点
        vector<Point> m_curve;  // 曲线上的点
        int m_k;                // B样条曲线次数
        vector<double> m_t;     // 节点向量
    };

    // 全局变量
    BSpline spline;
    bool dragging = false;

    // 初始化函数
    void init()
    {
        glClearColor(1.0, 1.0, 1.0, 0.0);
        glMatrixMode(GL_PROJECTION);
        glLoadIdentity();
        gluOrtho2D(0, 800, 0, 600);
    }

    // 显示函数
    void display()
    {
        glClear(GL_COLOR_BUFFER_BIT);
        spline.draw();
        glutSwapBuffers();
    }

    // 鼠标按下事件处理函数
    void mouse(int button, int state, int x, int y)
    {
        if (button == GLUT_LEFT_BUTTON && state == GLUT_DOWN)
        {
            double mx = x;
            double my = glutGet(GLUT_WINDOW_HEIGHT) - y;
            spline.addPoint(mx, my);
            spline.compute();
            glutPostRedisplay();
        }
    }

    // 鼠标移动事件处理函数
    void motion(int x, int y)
    {
        if (dragging)
        {
            double mx = x;
            double my = glutGet(GLUT_WINDOW_HEIGHT) - y;
            spline.addPoint(mx, my);
            spline.compute();
            glutPostRedisplay();
        }
    }

    // 鼠标拖拽事件处理函数
    void drag(int x, int y)
    {
        dragging = true;
        double mx = x;
        double my = glutGet(GLUT_WINDOW_HEIGHT) - y;
        spline.addPoint(mx, my);
        spline.compute();
        glutPostRedisplay();
    }

    // 鼠标松开事件处理函数
    void release(int x, int y)
    {
        dragging = false;
    }
}

// 主函数
int main(int argc, char **argv)
{
    glutInit(&argc, argv);
    glutInitDisplayMode(GLUT_DOUBLE | GLUT_RGB);
    glutInitWindowSize(800, 600);
    glutCreateWindow("B-Spline Curve");
    init();
    glutDisplayFunc(display);
    glutMouseFunc(mouse);
    glutMotionFunc(motion);
    glutPassiveMotionFunc(motion);
    glutKeyboardFunc(NULL);
    glutIdleFunc(NULL);
    glutDragFunc(drag);
    glutCloseFunc(NULL);
    glutEntryFunc(NULL);
    glutSpecialFunc(NULL);
    glutReleaseFunc(release);
    glutMainLoop();
    return 0;
}