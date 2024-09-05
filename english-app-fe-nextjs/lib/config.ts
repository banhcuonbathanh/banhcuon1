// const dotenv = require("dotenv");
// dotenv.config();

export const fetch_headers = {
  "Content-Type": "application/json"
};
export const HTTP_request_method_Get = {
  method: "GET",
  headers: fetch_headers
};
export const HTTP_request_method_Delete = {
  method: "DELETE",
  headers: fetch_headers
};
export const linkCustomer = {
  // golang_Base: process.env.golang_Base_Local,
  //  API_URL :  'http://127.0.0.1:8888',
  //  WEBSOCKET_URL : 'ws://127.0.0.1:8080',
  // golang_Base: process.env.NEXT_PUBLIC_DB_HOST || "http://shop-golang:8888",
  // WEBSOCKET_URL: "ws://shop-golang:8888",
  golang_Base: process.env.NEXT_PUBLIC_DB_HOST || "http://localhost:8888",
  WEBSOCKET_URL: "ws://localhost:8888",
  landingpagenavbarlogo: "/uploads/landingpage/navbar/logo.png",
  landingpagenavbarhumber: "/uploads/landingpage/navbar/hamburger.svg",
  shopnavbarlogo: "/uploads/shop/navbar/logo.png",
  landingpagenavbarstylistunserline:
    "/uploads/landingpage/herosection/stylish-underline.png",

  landingpagenavbarlady:
    "/uploads/landingpage/herosection/lady-with-aircraft.png",
  billboardhomeImageLink:
    "http://localhost:8888/uploads/b1/b1/home_bil_noard.png",

  API_URL: "http://localhost:8888",
  // auth
  loginImage: "/uploads/auth/login/login.png",
  signupImage: "/uploads/auth/signup/signup.png",
  // top selling section
  topsellingnavbarimage1: "/uploads/landingpage/topselling/europe.png",
  topsellingnavbarimage2: "/uploads/landingpage/topselling/rome.jpn",
  topsellingnavbarimage3: "/uploads/landingpage/topselling/london.jpg",
  topsellingsendicons: "/uploads/landingpage/topselling/send-icon.png",
  topsellingstylisring: "/uploads/landingpage/topselling/stylish-ring.png",
  //
  // golang_Base: "http://localhost:8888",
  routes_user: {
    getUsers: "/api/users",
    registerPostUser: "/api/authentication/register", // Renamed key

    getAuthenticationUsers: "/api/authentication/users",
    refreshToken: "/api/refresh-token",
    Login: "/api/users/login"
  },

  routes_image: {
    uploadImage: "/upload",
    serveImage: "/image",

    deleteImage: "//deleteImage/",
    serveThumbnail: "/thumbnail",
    serveStaticFile: "/uploads/*filepath"
  },
  routes_tags: {
    getAllTags: "/api/tags",
    getTagById: "/api/tags/",
    createTag: "/api/tags",
    updateTag: "/api/tags/",
    deleteTag: "/api/tags/"
  },
  routes_categories: {
    getAllCategories: "/api/categories",
    findByCategoryName: "/api/categories/findbyname",
    findCategoryById: "/api/categories/findbyid/",
    createCategory: "/api/categories",
    updateCategory: "/api/categories/",
    deleteCategory: "/api/categories/",
    findCategoriesByPage: "/api/categories/findbypage"
  },
  routes_size: {
    getAllSizes: "/api/sizes",
    getSizeById: "/api/sizes/",
    createSize: "/api/sizes",
    updateSize: "/api/sizes/",
    deleteSize: "/api/sizes/",
    findSizeByName: "/api/sizes/findbyname",
    findSizesByPage: "/api/sizes/findbypage"
  },
  routes_color: {
    getAllColors: "/api/colors",
    getColorById: "/api/colors/",
    createColor: "/api/colors",
    updateColor: "/api/colors/",
    deleteColor: "/api/colors/",
    findColorByName: "/api/colors/findbyname",
    findColorsByPage: "/api/colors/findbypage"
  },
  routes_billboards: {
    getAllBillboards: "/api/billboards",
    getBillboardById: "/api/billboards/",
    createBillboard: "/api/billboards",
    updateBillboard: "/api/billboards/",
    deleteBillboard: "/api/billboards/",
    findBillboardByLabel: "/api/billboards/findbylabel",
    findBillboardsByPage: "/api/billboards/findbypage"
  },
  routes_products: {
    getAllProducts: "/api/products",
    getProductById: "/api/products/findbyId/",
    createProduct: "/api/products",
    updateProduct: "/api/products/",
    deleteProduct: "/api/products/",
    findProductByCategory: "/api/products/findbycategory",
    findProductsByPriceRange: "/api/products/findbypricerange",
    findProductsByColor: "/api/products/findbycolor",
    findProductsBySize: "/api/products/findbysize",
    findProductsByCharacteristic: "/api/products/findbycharacteristic",
    findProductsByPage: "/api/products/findbypage",
    findProductsByName: "/api/products/findbyproductname"
  },

  routes_users: {
    createUser: "/api/users",
    updateUser: "/api/users/",
    deleteUser: "/api/users/",
    findUserById: "/api/users/",
    findAllUsers: "/api/users/findAll",
    findUserByEmail: "/api/users/findbyemail",
    findUsersByPage: "/api/users/findbypage"
  },
  routes_tokens: {
    createToken: "/api/tokens",
    getToken: "/api/tokens/:token",
    getTokensByUserID: "/api/tokens/user/:userID",
    updateToken: "/api/tokens",
    deleteToken: "/api/tokens/:token",
    getAllTokens: "/api/tokens"
  },

  routes_orders: {
    // Create order
    createOrder: {
      method: "POST",
      path: "/api/orders"
    },
    // # Find all orders
    findAllOrders: {
      method: "GET",
      path: "/api/orders"
    },
    // # Find order by ID
    findOrderById: {
      method: "GET",
      path: "/api/orders/:orderID"
    },
    // # Update order
    updateOrder: {
      method: "PATCH",
      path: "/api/orders/:orderID"
    },
    // # Delete order
    deleteOrder: {
      method: "DELETE",
      path: "/api/orders/"
    }
  },

  routes_orderItems: {
    createOrderItem: "/api/order-items",
    updateOrderItem: "/api/order-items/",
    deleteOrderItem: "/api/order-items/",

    getAllOrderItems: "/api/order-items",
    findOrderItemById: "/api/order-items/findbyId/",
    findOrderItemsByOrderID: "/api/order-items/findbyorderID/",
    findOrderItemsByPage: "/api/order-items/findbypage/",
    findOrderItemsByCustomerId: "/api/order-items/findbycustomer/"
  },

  routes_blogs: {
    // Get all blogs (potentially with pagination and filtering)
    findAllBlogs: {
      method: "GET",
      path: "/api/blogs",
      options: {
        pagination: {
          enabled: true,
          defaultPageSize: 10,
          maxPageSize: 100
        },
        filtering: {
          supportedFields: [
            "title",
            "category",
            "author",
            "createdAt",
            "updatedAt"
          ] // Example supported fields
        }
      }
    },
    // Get a specific blog by ID
    findBlogById: {
      method: "GET",
      path: "/api/blog-posts/:blogId"
    },
    // Create a new blog
    createBlog: {
      method: "POST",
      path: "/api/blog-posts"
    },
    // Update an existing blog
    updateBlog: {
      method: "PATCH",
      path: "/api/blog-posts/:blogId"
    },
    // Delete a blog
    deleteBlog: {
      method: "DELETE",
      path: "/api/blog-posts/:blogId"
    }
  },
  routes_blog_categories: {
    createBlogCategory: "/api/blog-categories",

    getAllBlogCategories: "/api/blog-categories",
    findByBlogCategoryName: "/api/blog-categories/findbyname",
    findBlogCategoryById: "api/blog-categories/",

    updateBlogCategory: "/api/blog-categories/",
    deleteBlogCategory: "/api/blog-categories/",
    findBlogCategoriesByPage: "/api/blog-categories/findbypage"
  },
  routes_blog_sub_des: {
    saveList: "/api/blog-sub-des/saveList",
    create: "/api/blog-sub-des",
    get: "/api/blog-sub-des/:subdesId",
    update: "/api/blog-sub-des/:subdesId",
    delete: "/api/blog-sub-des/:subdesId",
    getAll: "/api/blog-sub-des",
    search: "/api/blog-sub-des/search"
  },
  routes_blog_posts: {
    getAll: "/api/blog-posts",
    get: "/api/blog-posts/:postId",
    create: "/api/blog-posts",
    update: "/api/blog-posts/:postId",
    delete: "/api/blog-posts/:postId"
  },

  routes_blog_comments: {
    FindAllCommentsByPostIDFE: "/api/blog-comments",
    create: "/api/blog-comments",
    get: "/api/blog-comments/:commentID",
    update: "/api/blog-comments/:commentID",
    delete: "/api/blog-comments/:commentID",
    FindAllCommentsByPostID: "/api/blog-comments/FindAllCommentsByPostID/"
  },

  routes_chat: {
    createRoom: "/api/chat/createRoom",
    joinRoom: "/api/chat/joinRoom/",
    getRooms: "/api/chat/getRooms",
    getClients: "/api/chat/getClients/:roomId"
  },

  routes_blog_new_comments: {
    createComment: "/api/blog-new-comments",
    findAllCommentsByPostID: "/api/blog-new-comments/FindAllCommentsByPostID/]",
    updateComment: "/api/blog-new-comments/",
    deleteComment: "/api/blog-new-comments/",
    findCommentByID: "/api/blog-new-comments/",
    addBlogNewCommentModel: "/api/blog-new-comments/AddBlogNewCommentModel"
  },

  routes_ws: {
    createRoom: "/ws/createRoom",

    getRooms: "/ws/getRooms",
    joinRooms: "/ws/joinRoom/",
    getClients: "/ws/getClients/",
    WEBSOCKET_URL: "ws://localhost:8888"
  }
};

export const link_internal = {
  routes_shop_admin_categories: {
    categoryid: "/api/shop_admin/categories",
    addfress: "/admin/categories"
  },
  routes_shop_admin_categories_not_api: {
    categoryid: "/shop_admin/categories"
  },
  routes_shop_admin_sizes: {
    sizeid: "/api/shop_admin/sizes",

    addfress: "/admin/sizes"
  },
  routes_shop_admin_size_not_api: {
    sizeid: "/admin/sizes/"
  },
  routes_shop_admin_colors: {
    colorid: "/api/shop_admin/colors/"
  },
  routes_shop_admin_colors_not_api: {
    colorid: "/shop_admin/colors/"
  },
  routes_shop_admin_billboards: {
    billboardid: "/api/shop_admin/billboards/"
  },
  routes_shop_admin_billboards_not_api: {
    billboardid: "/shop_admin/billboards/"
  },
  routes_shop_admin_products: {
    productid: "/api/shop_admin/products/"
  },
  routes_shop_admin_products_not_api: {
    productid: "/shop_admin/products/"
  }
};
