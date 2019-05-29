package org.coolsugar.log.vo;

import java.io.Serializable;

//用户外部展示，字段可以涉及多个表的部分字段
public class OutUser implements Serializable {
    private Integer id;
    private Integer creator;
    private String name;

    public Integer getId() {
        return id;
    }

    public void setId(Integer id) {
        this.id = id;
    }

    public Integer getCreator() {
        return creator;
    }

    public void setCreator(Integer creator) {
        this.creator = creator;
    }

    public String getName() {
        return name;
    }

    public void setName(String name) {
        this.name = name == null ? null : name.trim();
    }

    @Override
    public String toString() {
        return "OutUser{" +
                "id=" + id +
                ", creator=" + creator +
                ", name='" + name +
                '}';
    }
}
